package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

func (server *Server) TestEncryptedUpload(c *gin.Context) {

	bucketname := "test-encrypted"
	err := utils.CreateBucketIfMissing(c.Request.Context(), *server.Minio, bucketname)
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("bucket creation failed")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// objectName := "my-objectname"

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("reading file failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fileName := c.Request.PostFormValue("name")
	if len(fileName) == 0 {
		logrus.WithField(
			"error", err,
		).Error("reading file name failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	size, err := strconv.ParseInt(c.Request.PostFormValue("size"), 10, 64)
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("reading file size failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logrus.WithFields(logrus.Fields{
		"file": file,
	}).Info("file request")

	// in app we can have something like lock
	// each user will have encrypted and non encrypted buckets
	password := "abc"
	// New SSE-C where the cryptographic key is derived from a password and the objectname + bucketname as salt
	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+fileName))
	// Encrypt file content and upload to the server
	uploadedInfo, err := server.Minio.PutObject(context.Background(), bucketname, fileName, file, int64(size), minio.PutObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("file upload failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logrus.WithField(
		"fileInfo", uploadedInfo,
	).Info("file uploaded")
}

func (server *Server) TestNormalUpload(c *gin.Context) {
	bucketname := "test"
	err := utils.CreateBucketIfMissing(c.Request.Context(), *server.Minio, bucketname)
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("bucket creation failed")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// objectName := "my-objectname"

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("reading file failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fileName := c.Request.PostFormValue("name")
	if len(fileName) == 0 {
		logrus.WithField(
			"error", err,
		).Error("reading file name failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	size, err := strconv.ParseInt(c.Request.PostFormValue("size"), 10, 64)
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("reading file size failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logrus.WithFields(logrus.Fields{
		"file":   file,
		"header": header,
	}).Info("file request")

	// Encrypt file content and upload to the server
	uploadedInfo, err := server.Minio.PutObject(context.Background(), bucketname, fileName, file, int64(size), minio.PutObjectOptions{})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logrus.WithField(
		"fileInfo", uploadedInfo,
	).Info("file uploaded")
}

func (server *Server) TestGetEncryptedVideo(c *gin.Context) {
	rangeHeader := c.Request.Header["Range"]
	var parsedRangeHeader *RangeHeader
	if len(rangeHeader) != 0 && len(rangeHeader[0]) != 0 {
		var err error
		parsedRangeHeader, err = parseRangeHeader(rangeHeader[0])
		if err != nil {
			logrus.Error(err)
		}
	}
	fileName := c.Request.FormValue("file")
	if len(fileName) == 0 {
		logrus.WithField(
			"file", fileName,
		).Error("not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	bucketname := "test-encrypted"
	password := "abc"
	// New SSE-C where the cryptographic key is derived from a password and the objectname + bucketname as salt
	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+fileName))
	// todo no error on object not existing
	object, err := server.Minio.GetObject(c.Request.Context(), bucketname, fileName, minio.GetObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if parsedRangeHeader == nil || len(parsedRangeHeader.ranges) != 1 {
		server.TestGetVideoFirstRequest(c, object)
		return
	}
	server.TestGetRangeVideo(c, parsedRangeHeader.ranges[0], object)
}

// func (server *Server) GetEncryptedVideo(c *gin.Context) {
// 	bucketname := "test-encrypted"
// 	fileName := "test.mp4"
// 	// in app we can have something like lock
// 	// each user will have encrypted and non encrypted buckets
// 	password := "abc"
// 	// New SSE-C where the cryptographic key is derived from a password and the objectname + bucketname as salt
// 	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+fileName))
// 	// Encrypt file content and upload to the server
// 	object, err := server.Minio.GetObject(context.Background(), bucketname, fileName, minio.GetObjectOptions{ServerSideEncryption: encryption})

// 	// url, err := server.Minio.PresignedGetObject(context.Background(), bucketname, fileName, minio.GetObjectOptions{ServerSideEncryption: encryption})

// }
