package v1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

// seek reader?
// todo add mutex
type uploadRequest struct {
	Reader     io.ReadCloser
	Writer     io.WriteCloser
	err        error // in case of failure
	completed  bool
	index      int64
	size       int64
	ctx        context.Context
	cancelFunc context.CancelFunc
	checksum   string
}

// todo session affinity required till all the browsers support http/2 protocol (which support stream upload)
// https://caniuse.com/http2, right now android browser's don't have good support
// requestID -> uploadRequest
var uploadRequests map[string]*uploadRequest = make(map[string]*uploadRequest)

func newUploadRequest(ctx context.Context, cancelFunc context.CancelFunc, requestID string) error {
	if _, ok := uploadRequests[requestID]; ok {
		return fmt.Errorf("request with ID %v already exist", requestID)
	}
	reader, writer := io.Pipe()
	uploadRequests[requestID] = &uploadRequest{
		Reader:     reader,
		Writer:     writer,
		err:        nil,
		completed:  false,
		index:      0,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
	return nil
}

func deleteUploadRequestAfter(t time.Duration, requestID string) error {
	if _, ok := uploadRequests[requestID]; !ok {
		return fmt.Errorf("request with ID %v doesn't Exist", requestID)
	}
	// wait for the time and then delete the key from map
	<-time.NewTicker(t).C
	delete(uploadRequests, requestID)
	return nil
}

type initRequestBody struct {
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
	FileType string `json:"fileType"`
}

func (server *Server) InitChunkUpload(c *gin.Context) {
	var requestBody initRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logrus.Infof("[InitUpload] invalid request: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	requestID := uuid.New().String()
	go server.startUploadInBackground(requestID, requestBody.FileName, requestBody.Size)
	c.JSON(http.StatusOK, gin.H{
		"requestID": requestID,
	})
}

func (server *Server) startUploadInBackground(requestID string, fileName string, size int64) {
	// todo upgrade go and change this to WithCancelCause
	ctx, cancelFunc := context.WithCancel(context.Background())
	err := newUploadRequest(ctx, cancelFunc, requestID)
	if err != nil {
		logrus.Errorf("[server.startUploadInBackground] request creation failed: %v", err)
		return
	}
	bucketName := "test"
	uploadRequest, ok := uploadRequests[requestID]
	if !ok {
		logrus.Errorf("[server.startUploadInBackground] request %v not found", requestID)
		return
	}
	err = utils.CreateBucketIfMissing(ctx, *server.Minio, bucketName)
	if err != nil {
		logrus.Errorf("[server.startUploadInBackground] bucket creation failed: %v", err)
		return
	}
	uploadRequest.size = size
	// todo need to add some kind of timeout during upload if no data is transfered for sometime
	// i think tcp by default has some timeout
	uploadedInfo, err := server.Minio.PutObject(ctx, bucketName, fileName, uploadRequest.Reader, size, minio.PutObjectOptions{})
	if err != nil {
		// todo time="2023-02-19T09:22:39Z" level=error msg="[server.startUploadInBackground] upload failed: A timeout occurred while trying to lock a resource, please reduce your request rate"
		logrus.Errorf("[server.startUploadInBackground] upload failed: %v", err)
		uploadRequest.completed = true
		uploadRequest.err = err
		uploadRequest.cancelFunc()
		// delete the request after 10 minutes to free memory
		// finishUpload request will not work after 10 minutes, so client only has 10 minutes to send finishUpload request
		// todo call delete right away and store rest of the details in DB instead
		deleteUploadRequestAfter(1*time.Hour, requestID)
		return
	}
	logrus.Infof("[server.startUploadInBackground] upload completed: %v", uploadedInfo)
	uploadRequest.completed = true
	uploadRequest.err = nil
	uploadRequest.cancelFunc()
	// delete the request after 10 minutes to free memory
	// finishUpload request will not work after 10 minutes, so client has 10 minutes
	deleteUploadRequestAfter(10*time.Minute, requestID)
}

// type uploadChunkRequestBody struct {
// 	RequestID string `json:"requestID"`
// 	ChunkSize int64  `json:"chunkSize"`
// 	// type?
// 	ChunkData io.Reader `json:"chunkData"`
// 	Index     int64     `json:"index"`
// }

// multipar request
func (server *Server) UploadChunk(c *gin.Context) {
	requestID := c.Request.PostFormValue("requestID")
	if len(requestID) == 0 {
		logrus.Error("[Server.uploadChunk] bad request: requestID param missing")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var err error
	var index int64
	{
		value := c.Request.PostFormValue("index")
		if len(value) == 0 {
			logrus.Error("[Server.uploadChunk] bad request: index param missing")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		index, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			logrus.Error("[Server.uploadChunk] bad request: invalid index %v", value)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	var chunkSize int64
	{
		value := c.Request.PostFormValue("chunkSize")
		if len(value) == 0 {
			logrus.Error("[Server.uploadChunk] bad request: chunkSize param missing")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		chunkSize, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			logrus.Error("[Server.uploadChunk] bad request: invalid chunkSize %v", value)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	chunkData, _, err := c.Request.FormFile("chunkData")
	if err != nil {
		logrus.Error("[Server.uploadChunk] bad request: chunkData param missing")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	uploadRequest, ok := uploadRequests[requestID]
	if !ok {
		logrus.Errorf("[Server.uploadChunk] bad request: no uploadRequest found with the requestID %v", requestID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if uploadRequest.completed {
		logrus.Errorf("[Server.uploadChunk] bad request: upload is already completed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if index != uploadRequest.index {
		logrus.Errorf("[Server.uploadChunk] bad request: index mismatch expected %v but got %v", uploadRequest.index, index)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// todo buffer, it causes high cpu if we are sending a lot of small requests
	n, err := io.CopyN(uploadRequest.Writer, chunkData, chunkSize)
	if err != nil {
		logrus.Errorf("[Server.uploadChunk] io.CopyN failed: %v", err)
		uploadRequest.cancelFunc()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	uploadRequest.index += n
	if uploadRequest.index == uploadRequest.size {
		if err := uploadRequest.Writer.Close(); err != nil {
			logrus.Errorf("[server.uploadChunk] error closing writer: %v", err)
		}
		// if err := uploadRequest.Reader.Close(); err != nil {
		// 	logrus.Errorf("[server.uploadChunk] error closing reader: %v", err)
		// }
	}
	c.JSON(http.StatusOK, gin.H{
		"requestID": requestID,
		"uploaded":  n,
	})
}

type finishUploadRequestBody struct {
	RequestID string `json:"requestID"`
	Checksum  string `json:"checksum"`
}

func (server *Server) FinishChunkUpload(c *gin.Context) {
	var requestBody finishUploadRequestBody
	c.Request.ParseForm()
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		logrus.Errorf("[server.FinishUpload] bad request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	uploadRequest, ok := uploadRequests[requestBody.RequestID]
	if !ok {
		logrus.Errorf("[server.FinishUpload] bad request: request %v does not exist", requestBody.RequestID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !uploadRequest.completed {
		logrus.Infof("[server.FinishUpload]: upload request %v should have completed but it is still not completed yet. last chuck upload to minio might be still in progress will wait for 1 more minute", requestBody.RequestID)

		waitTime := 1 * time.Minute
		select {
		case <-time.NewTicker(waitTime).C:
			logrus.Errorf("[server.FinishUpload]: upload did not complete in time")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		case <-uploadRequest.ctx.Done():
			break
		}
	}
	if uploadRequest.err != nil {
		logrus.Errorf("[server.FinishUpload] upload failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	return
}

// todo if upload fail delete file in minio
