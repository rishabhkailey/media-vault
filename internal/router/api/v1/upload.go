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
