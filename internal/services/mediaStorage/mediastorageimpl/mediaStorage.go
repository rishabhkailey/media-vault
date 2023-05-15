package mediastorageimpl

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	mediastorage "github.com/rishabhkailey/media-service/internal/services/mediaStorage"
	uploadrequests "github.com/rishabhkailey/media-service/internal/services/uploadRequests"
	"github.com/sirupsen/logrus"
)

type Service struct {
	store                 store
	onGoingUploadRequests onGoingUploadRequestsStore
	uploadRequestsService uploadrequests.Service
}

var _ mediastorage.Service = (*Service)(nil)

func NewMinioService(cli *minio.Client, bucketName string, uploadRequestsService uploadrequests.Service,
) (*Service, error) {
	var store store
	store, err := newMinioStore(cli, bucketName)
	if err != nil {
		return nil, err
	}
	store, err = NewStoreCacheWrapper(store)
	if err != nil {
		return nil, err
	}
	return &Service{
		store:                 store,
		onGoingUploadRequests: newOnGoingUploadRequestsStore(),
		uploadRequestsService: uploadRequestsService,
	}, nil
}

func (s *Service) GetMediaByFileName(ctx context.Context, query mediastorage.GetMediaByFileNameQuery) (mediastorage.File, error) {
	return s.store.GetByFileName(ctx, query.FileName)
}

func (s *Service) GetThumbnailByFileName(ctx context.Context, query mediastorage.GetThumbnailByFileNameQuery) (mediastorage.File, error) {
	return s.store.GetByFileName(ctx, s.GetThumbnailFileName(query.FileName))
}

func (s *Service) HttpGetRangeHandler(ctx context.Context, query mediastorage.HttpGetRangeHandlerQuery) (int64, error) {
	file, err := s.store.GetByFileName(ctx, query.FileName)
	if err != nil {
		return 0, fmt.Errorf("[mediaService.HttpGetRangeHandler]: get by file name failed: %w", err)
	}
	stat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("[mediaService.HttpGetRangeHandler]: get file stats failed: %w", err)
	}
	// range end inclusive
	if query.Range.End > stat.Size()-1 {
		query.Range.End = stat.Size() - 1
	}
	_, err = file.Seek(query.Range.Start, 0)
	if err != nil {
		return 0, fmt.Errorf("[mediaService.HttpGetRangeHandler]: file seek failed: %w", err)
	}
	// range start and end is inclusive
	contentLength := query.Range.End - query.Range.Start + 1
	// todo why can't we move the header logic to router and here we only copy the data
	// change the function names accordingly
	query.ResponseWriter.Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", query.Range.Start, query.Range.End, stat.Size()))
	// query.ResponseWriter.Header().Add("content-length", strconv.FormatInt(contentLength, 10))
	return io.CopyN(query.ResponseWriter, file, contentLength)
}

func (s *Service) HttpGetMediaHandler(ctx context.Context, query mediastorage.HttpGetMediaHandlerQuery) (int64, error) {
	file, err := s.store.GetByFileName(ctx, query.FileName)
	if err != nil {
		return 0, err
	}
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	// query.ResponseWriter.Header().Add("content-length", strconv.FormatInt(stat.Size(), 10))
	return io.CopyN(query.ResponseWriter, file, stat.Size())
}

func (s *Service) HttpGetThumbnailHandler(ctx context.Context, query mediastorage.HttpGetThumbnailHandlerQuery) (int64, error) {
	return s.HttpGetMediaHandler(ctx, mediastorage.HttpGetMediaHandlerQuery{
		FileName:       s.GetThumbnailFileName(query.FileName),
		ResponseWriter: query.ResponseWriter,
	})
}

func (s *Service) InitChunkUpload(_ context.Context, cmd mediastorage.InitChunkUploadCmd) error {
	// todo upgrade go and change this to WithCancelCause
	uploadCtx, cancelFunc := context.WithCancel(context.Background())
	uploadRequest, err := s.onGoingUploadRequests.add(uploadCtx, cancelFunc, cmd)
	if err != nil {
		return err
	}
	go func(ctx context.Context) {
		// todo need to add some kind of timeout during upload if no data is transfered for sometime
		// i think tcp by default has some timeout
		n, err := s.store.SaveFile(ctx, mediastorage.StoreSaveFileCmd{
			FileName:   cmd.FileName,
			FileSize:   cmd.FileSize,
			FileReader: uploadRequest.Reader,
		})
		if err != nil {
			// todo time="2023-02-19T09:22:39Z" level=error msg="[server.startUploadInBackground] upload failed: A timeout occurred while trying to lock a resource, please reduce your request rate"
			logrus.Errorf("[server.startUploadInBackground] upload failed: %v", err)
			uploadRequest.completed = true
			uploadRequest.err = err
			uploadRequest.cancelFunc()
			s.onGoingUploadRequests.deleteUploadRequestAfter(0*time.Second, cmd.RequestID, cmd.UserID)
			// todo services should be loosely coupled
			// try to move this to onFailure function?
			err := s.uploadRequestsService.UpdateStatus(context.Background(), uploadrequests.UpdateStatusCommand{
				ID:     cmd.RequestID,
				Status: uploadrequests.FAILED_UPLOAD_STATUS,
			})
			if err != nil {
				logrus.Errorf("[server.startUploadInBackground] uploadRequest update status failed: %v", err)
			}
			return
		}
		logrus.Infof("[server.startUploadInBackground] upload completed: %n bytes", n)
		uploadRequest.completed = true
		uploadRequest.err = nil
		uploadRequest.cancelFunc()
		// todo services should be loosely coupled
		// try to move this to onSuccess function?
		err = s.uploadRequestsService.UpdateStatus(context.Background(), uploadrequests.UpdateStatusCommand{
			ID:     cmd.RequestID,
			Status: uploadrequests.COMPLETED_UPLOAD_STATUS,
		})
		if err != nil {
			logrus.Errorf("[server.startUploadInBackground] uploadRequest update status failed: %v", err)
		}
		// delete the request after 10 minutes to free memory
		// finishUpload request will not work after 10 minutes, so client has 10 minutes
		s.onGoingUploadRequests.deleteUploadRequestAfter(10*time.Minute, cmd.RequestID, cmd.UserID)

	}(uploadCtx)
	return nil
}

func (s *Service) UploadChunk(ctx context.Context, cmd mediastorage.UploadChunkCmd) (int64, error) {
	uploadRequest, err := s.onGoingUploadRequests.get(cmd.UploadRequestID, cmd.UserID)
	if err != nil {
		return 0, fmt.Errorf("[storageService.UploadChunk] get onGoingUploadRequest failed: %w", err)
	}
	if cmd.UserID != uploadRequest.userID {
		return 0, fmt.Errorf("[storageService.UploadChunk] unauthorized: incorrest user")
	}
	if uploadRequest.completed {
		return 0, fmt.Errorf("[storageService.UploadChunk] uploadRequest %s is already completed. Bad request", cmd.UploadRequestID)
	}
	if cmd.Index != uploadRequest.index {
		return 0, fmt.Errorf("[storageService.UploadChunk] index mismatch, possible bad request. expected %d got %d", uploadRequest.index, cmd.Index)
	}
	n, err := io.CopyN(uploadRequest.Writer, cmd.Chunk, cmd.ChunkSize)
	if err != nil {
		uploadRequest.cancelFunc()
		return n, fmt.Errorf("[storageService.UploadChunk] upload chunk failed: %w", err)
	}
	uploadRequest.index += n
	if uploadRequest.index == uploadRequest.size {
		if err := uploadRequest.Writer.Close(); err != nil {
			logrus.Errorf("[server.uploadChunk] error closing writer, possible memroy leak: %v", err)
		}
	}
	return n, nil
}

func (s *Service) FinishChunkUpload(ctx context.Context, cmd mediastorage.FinishChunkUpload) error {
	uploadRequest, err := s.onGoingUploadRequests.get(cmd.RequestID, cmd.UserID)
	if err != nil {
		return fmt.Errorf("[storageService.FinishChunkUpload] bad request: request %v does not exist. %w", cmd.RequestID, err)
	}
	if cmd.UserID != uploadRequest.userID {
		return fmt.Errorf("[storageService.FinishChunkUpload] unauthorized: incorrest user")
	}
	if !uploadRequest.completed {
		logrus.Infof("[storageService.FinishChunkUpload]: upload request %v should have completed but it is still not completed yet. last chuck upload to minio might be still in progress will wait for 5 more minute", cmd.RequestID)
		waitTime := 5 * time.Minute
		select {
		case <-time.NewTicker(waitTime).C:
			return fmt.Errorf("[storageService.FinishChunkUpload]: upload did not complete in time")
		case <-uploadRequest.ctx.Done():
			break
		}
	}
	if uploadRequest.err != nil {
		return fmt.Errorf("[storageService.FinishChunkUpload] upload failed: %v", err)
	}
	return nil
}

func (s *Service) ThumbnailUpload(ctx context.Context, cmd mediastorage.UploadThumbnailCmd) error {
	uploadRequest, err := s.onGoingUploadRequests.get(cmd.RequestID, cmd.UserID)
	if err != nil {
		return fmt.Errorf("[storageService.ThumbnailUpload] bad request: request %v does not exist. %w", cmd.RequestID, err)
	}
	if cmd.UserID != uploadRequest.userID {
		return fmt.Errorf("[storageService.ThumbnailUpload] unauthorized: incorrest user")
	}
	thumbnailFileName := s.GetThumbnailFileName(cmd.FileName)
	_, err = s.store.SaveFile(ctx, mediastorage.StoreSaveFileCmd{
		FileName:   thumbnailFileName,
		FileSize:   cmd.FileSize,
		FileReader: cmd.FileReader,
	})
	return err
}

func (s *Service) GetThumbnailFileName(fileName string) string {
	return fmt.Sprintf(".thumb-%s", fileName)
}
