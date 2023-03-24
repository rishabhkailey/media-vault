package v1

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session/v3"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/auth"
	dbservices "github.com/rishabhkailey/media-service/internal/db/services"
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
	// checksum   string
	// mediaID    uint
	// userID     uint
}

// todo session affinity required till all the browsers support http/2 protocol (which support stream upload)
// https://caniuse.com/http2, right now android browser's don't have good support
// requestID -> uploadRequest
var uploadRequests map[string]*uploadRequest = make(map[string]*uploadRequest)

func newUploadRequest(ctx context.Context, cancelFunc context.CancelFunc, requestID, userID string) error {
	if _, ok := uploadRequests[requestID]; ok {
		return fmt.Errorf("request with ID %v already exist", requestID)
	}
	reader, writer := io.Pipe()
	uploadRequests[fmt.Sprintf("%s:%s", requestID, userID)] = &uploadRequest{
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

func getUploadRequest(requestID, userID string) (*uploadRequest, error) {
	uploadRequest, ok := uploadRequests[fmt.Sprintf("%s:%s", requestID, userID)]
	if !ok {
		return nil, fmt.Errorf("request with ID %v doesn't Exist", requestID)
	}
	return uploadRequest, nil
}

func deleteUploadRequestAfter(t time.Duration, requestID, userID string) error {
	key := fmt.Sprintf("%s:%s", requestID, userID)
	if _, ok := uploadRequests[key]; !ok {
		return fmt.Errorf("request with ID %v doesn't Exist", requestID)
	}
	<-time.NewTicker(t).C
	delete(uploadRequests, key)
	return nil
}

type initRequestBody struct {
	FileName  string `json:"fileName"`
	Size      int64  `json:"size"`
	MediaType string `json:"mediaType,omitempty"`
	Date      int64  `json:"date,omitempty"`
}

func (server *Server) InitChunkUpload(c *gin.Context) {
	token, ok := auth.GetBearerToken(c.Request)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenInfo, err := server.OidcClient.IntrospectToken(token)
	if errors.Is(err, auth.ErrUnauthorized) || !tokenInfo.ValidateScope(auth.SCOPE_USER) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if len(tokenInfo.Subject) == 0 {
		logrus.Errorf("token info doesn't contain user info")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var requestBody initRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logrus.Infof("[InitUpload] invalid request: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	if len(requestBody.FileName) == 0 || requestBody.Size == 0 {
		logrus.Infof("[InitUpload] invalid request: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	if len(requestBody.MediaType) == 0 {
		requestBody.MediaType = string(dbservices.UNKNOWN)
	}
	if requestBody.Date == 0 {
		requestBody.Date = time.Now().Unix()
	}
	uploadRequest, err := server.UploadRequests.Create(c.Request.Context(), tokenInfo.Subject)
	if err != nil {
		logrus.Errorf("[InitUpload] uploadRequest creation failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	media, err := server.Media.Create(c.Request.Context(), *uploadRequest)
	if err != nil {
		logrus.Errorf("[InitUpload] media creation failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = server.MediaMetadata.Create(c.Request.Context(), *media, dbservices.Metadata{
		Name: requestBody.FileName,
		Date: time.Unix(requestBody.Date, 0),
		Size: uint64(requestBody.Size),
		Type: requestBody.MediaType,
	})
	if err != nil {
		logrus.Errorf("[InitUpload] media metadata creation failed: %w", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// for 1 upload request we will only validate token on init request and store info in the session
	// if session contains requestID:user = UserID then user is authenticated
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.InitUpload"}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}
	// todo - requestID to check if the request belongs to user or not
	// further userID will be used in the upload requests slice
	store.Set(fmt.Sprintf("%s:user", uploadRequest.ID), tokenInfo.Subject)
	store.Save()
	go server.startUploadInBackground(uploadRequest.ID, tokenInfo.Subject, media.FileName, requestBody.Size)
	c.JSON(http.StatusOK, gin.H{
		"requestID": uploadRequest.ID,
	})
}

func (server *Server) startUploadInBackground(requestID, userID, fileNameOnServer string, size int64) {
	// todo upgrade go and change this to WithCancelCause
	ctx, cancelFunc := context.WithCancel(context.Background())
	err := newUploadRequest(ctx, cancelFunc, requestID, userID)
	if err != nil {
		logrus.Errorf("[server.startUploadInBackground] request creation failed: %v", err)
		return
	}
	bucketName := "test"
	uploadRequest, err := getUploadRequest(requestID, userID)
	// how to inform user about these errors?
	if err != nil {
		logrus.Errorf("[server.startUploadInBackground] failed to get upload request ID: %w", err)
	}
	err = utils.CreateBucketIfMissing(ctx, *server.Minio, bucketName)
	if err != nil {
		logrus.Errorf("[server.startUploadInBackground] bucket creation failed: %v", err)
		return
	}
	uploadRequest.size = size
	// todo need to add some kind of timeout during upload if no data is transfered for sometime
	// i think tcp by default has some timeout
	uploadedInfo, err := server.Minio.PutObject(ctx, bucketName, fileNameOnServer, uploadRequest.Reader, size, minio.PutObjectOptions{})
	if err != nil {
		// todo time="2023-02-19T09:22:39Z" level=error msg="[server.startUploadInBackground] upload failed: A timeout occurred while trying to lock a resource, please reduce your request rate"
		logrus.Errorf("[server.startUploadInBackground] upload failed: %v", err)
		uploadRequest.completed = true
		uploadRequest.err = err
		uploadRequest.cancelFunc()
		deleteUploadRequestAfter(0*time.Second, requestID, userID)
		err := server.UploadRequests.UpdateStatus(context.Background(), requestID, dbservices.FAILED_UPLOAD_STATUS)
		if err != nil {
			logrus.Errorf("[server.startUploadInBackground] uploadRequest update status failed: %v", err)
		}
		return
	}
	logrus.Infof("[server.startUploadInBackground] upload completed: %v", uploadedInfo)
	uploadRequest.completed = true
	uploadRequest.err = nil
	uploadRequest.cancelFunc()
	err = server.UploadRequests.UpdateStatus(context.Background(), requestID, dbservices.COMPLETED_UPLOAD_STATUS)
	if err != nil {
		logrus.Errorf("[server.startUploadInBackground] uploadRequest update status failed: %v", err)
	}
	// delete the request after 10 minutes to free memory
	// finishUpload request will not work after 10 minutes, so client has 10 minutes
	deleteUploadRequestAfter(10*time.Minute, requestID, userID)
}

// type uploadChunkRequestBody struct {
// 	RequestID string `json:"requestID"`
// 	ChunkSize int64  `json:"chunkSize"`
// 	// type?
// 	ChunkData io.Reader `json:"chunkData"`
// 	Index     int64     `json:"index"`
// }

// todo check for multipart request
// todo request size limit
func (server *Server) UploadChunk(c *gin.Context) {
	requestID := c.Request.PostFormValue("requestID")
	if len(requestID) == 0 {
		logrus.Error("[Server.uploadChunk] bad request: requestID param missing")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.InitUpload"}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}
	var userID string
	if value, ok := store.Get(fmt.Sprintf("%s:user", requestID)); ok {
		userID, _ = value.(string)
	}
	if len(userID) == 0 {
		logrus.Errorf("[UploadChunk]: requestID not found in session")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
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
	uploadRequest, err := getUploadRequest(requestID, userID)
	if err != nil {
		logrus.Errorf("[server.uploadChunk] failed to get upload request ID: %w", err)
		// todo bad request?
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

// this is for client to confirm if the upload has finished or not
func (server *Server) FinishChunkUpload(c *gin.Context) {
	var requestBody finishUploadRequestBody
	c.Request.ParseForm()
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		logrus.Errorf("[server.FinishUpload] bad request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "function": "server.InitUpload"}).Errorf("session start failed")
		// todo error response
		c.Status(http.StatusInternalServerError)
		return
	}
	var userID string
	if value, ok := store.Get(fmt.Sprintf("%s:user", requestBody.RequestID)); ok {
		userID, _ = value.(string)
	}
	if len(userID) == 0 {
		logrus.Errorf("[UploadChunk]: requestID not found in session")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	uploadRequest, err := getUploadRequest(requestBody.RequestID, userID)
	if err != nil {
		logrus.Errorf("[server.FinishUpload] bad request: request %v does not exist", requestBody.RequestID)
		// todo bad request?
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !uploadRequest.completed {
		logrus.Infof("[server.FinishUpload]: upload request %v should have completed but it is still not completed yet. last chuck upload to minio might be still in progress will wait for 5 more minute", requestBody.RequestID)
		waitTime := 5 * time.Minute
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
}

// todo if upload fail delete file in minio
