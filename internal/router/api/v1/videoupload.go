package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var ErrIncorrectFileSize = errors.New("incorrect file size")

func (server *Server) TestVideoUploadWithThumbnail(c *gin.Context) {

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
	thumbnailName := fmt.Sprintf(".thumb-%s.%s", fileName, ".png")
	{
		name := strings.Split(fileName, ".")
		if len(name) > 1 {
			thumbnailName = fmt.Sprintf(".thumb-%s.%s", name[0], ".png")
		}
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

	ffmpegReader, ffmpegWriter := io.Pipe()
	minioReader, minioWriter := io.Pipe()
	commonWriter := io.MultiWriter(ffmpegWriter, minioWriter)

	// todo check error groups instead of using channel if that makes sense. as we don't have any data in done channel
	minioUploadErrs := make(chan error)
	thumbnailGeneratorErrs := make(chan error)
	minioUploadDone := make(chan bool)
	thumbnailGeneratorDone := make(chan bool)

	go server.uploadFileToMinio(c.Request.Context(), bucketname, fileName, minioReader, int64(size), minioUploadDone, minioUploadErrs)
	go server.generateAndUploadThumbnail(c.Request.Context(), bucketname, thumbnailName, ffmpegReader, size, thumbnailGeneratorDone, thumbnailGeneratorErrs)

	g, _ := errgroup.WithContext(c.Request.Context())
	g.Go(func() error {
		// todo if no progress in io copy then cancel the request
		defer ffmpegWriter.Close()
		defer minioWriter.Close()
		n, err := io.CopyN(commonWriter, file, size)
		if errors.Is(err, io.EOF) && n != size {
			logrus.Errorf("partially copied %v bytes", n)
			return ErrIncorrectFileSize
		}
		if err != nil {
			logrus.Infof("partially copied %v bytes", n)
			return err
		}
		logrus.Infof("successfuly copied %v bytes", n)
		return nil
	})
	err = g.Wait()
	if errors.Is(err, ErrIncorrectFileSize) {
		logrus.Errorf("Invalid request: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	if err != nil {
		logrus.Errorf("io.copy failed: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	// wait for minioUpload
	select {
	case <-minioUploadDone:
		logrus.Info("file uploaded to minio")
	case err := <-minioUploadErrs:
		logrus.Errorf("file upload failed: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// wait for thumbnail
	select {
	case <-thumbnailGeneratorDone:
		logrus.Info("thumbnail generated")
	case err := <-thumbnailGeneratorErrs:
		logrus.Errorf("file upload failed: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	return
}

func (server *Server) uploadFileToMinio(ctx context.Context, bucketname string, fileName string, file io.Reader, size int64, done chan<- bool, errs chan<- error) {
	// as we are using is multi writer it will block writer to minio writer the other.
	// defer io.Copy(io.Discard, file)
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				errs <- err
				return
			}
			errs <- fmt.Errorf("panic occured while generating thumbnails")
			return
		}
	}()
	uploadedInfo, err := server.Minio.PutObject(ctx, bucketname, fileName, file, size, minio.PutObjectOptions{})
	if err != nil {
		errs <- err
		return
	}
	logrus.WithField(
		"fileInfo", uploadedInfo,
	).Info("file uploaded")
	done <- true
}

// just rename to upload thumbnail + similar params as uploadFileToMinio
func (server *Server) generateAndUploadThumbnail(ctx context.Context, bucketname string, fileName string, file io.Reader, fileSize int64, done chan<- bool, errs chan<- error) {
	logrus.Info("generateThumbnail called")
	// as we are using is multi writer it will block writer to minio writer the other.
	// in case of failure we will need to do this else request will be stuck
	// todo better solutino for this
	defer io.Copy(io.Discard, file)
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				errs <- err
				return
			}
			errs <- fmt.Errorf("panic occured while generating thumbnails")
			return
		}
	}()
	thumbnailBytes, err := utils.GenerateThumbnail(file)
	if err != nil {
		errs <- err
		return
	}
	tempSaveThumbnail(thumbnailBytes)

	bytes.NewReader(thumbnailBytes)
	uploadedInfo, err := server.Minio.PutObject(ctx, bucketname, fileName, bytes.NewReader(thumbnailBytes), int64(len(thumbnailBytes)), minio.PutObjectOptions{})
	if err != nil {
		errs <- err
		return
	}
	logrus.WithField(
		"fileInfo", uploadedInfo,
	).Info("file uploaded")

	// todo check why putting done channel above this log doesn't work
	logrus.Infof("discarded ")
	logrus.Info(io.Copy(io.Discard, file))

	done <- true
}

func tempSaveThumbnail(b []byte) {
	file, err := os.Create("tmp.png")
	if err != nil {
		logrus.Error("create tmp.png failed")
		return
	}
	if _, err := io.Copy(file, bytes.NewReader(b)); err != nil {
		logrus.Error("tmp.png write failed")
	}
}

// func asyncIoCopy(w io.Writer, r io.Reader) error {
// 	n, err := io.Copy(w, r)
// 	if err != nil {
// 		return err
// 	}
// 	logrus.Infof("io.copy %v bytes", n)
// 	return nil
// }

// func (server *Server) TestParallelVideoUploadWithThumbnail(c *gin.Context) {

// 	bucketname := "test"
// 	err := utils.CreateBucketIfMissing(c.Request.Context(), *server.Minio, bucketname)
// 	if err != nil {
// 		logrus.WithField(
// 			"error", err,
// 		).Error("bucket creation failed")
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}
// 	// objectName := "my-objectname"

// 	file, _, err := c.Request.FormFile("file")
// 	if err != nil {
// 		logrus.WithField(
// 			"error", err,
// 		).Error("reading file failed")
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	fileName := c.Request.PostFormValue("name")
// 	if len(fileName) == 0 {
// 		logrus.WithField(
// 			"error", err,
// 		).Error("reading file name failed")
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	size, err := strconv.ParseInt(c.Request.PostFormValue("size"), 10, 64)
// 	if err != nil {
// 		logrus.WithField(
// 			"error", err,
// 		).Error("reading file size failed")
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	logrus.WithFields(logrus.Fields{
// 		"file": file,
// 	}).Info("file request")

// 	ffmpegReader, ffmpegWriter := io.Pipe()
// 	minioReader, minioWriter := io.Pipe()
// 	commonWriter := io.MultiWriter(ffmpegWriter, minioWriter)

// 	// todo check error groups instead of using channel if that makes sense. as we don't have any data in done channel
// 	minioUploadErrs := make(chan error)
// 	thumbnailGeneratorErrs := make(chan error)
// 	minioUploadDone := make(chan bool)
// 	thumbnailGeneratorDone := make(chan bool)

// 	go server.uploadFileToMinio(c.Request.Context(), bucketname, fileName, minioReader, int64(size), minioUploadDone, minioUploadErrs)
// 	// go func() {
// 	// 	logrus.Info(io.CopyN(io.Discard, ffmpegReader, size))
// 	// 	logrus.Info(size)
// 	// }()
// 	go server.generateThumbnail(c.Request.Context(), ffmpegReader, thumbnailGeneratorDone, thumbnailGeneratorErrs)

// 	g, _ := errgroup.WithContext(c.Request.Context())
// 	g.Go(func() error {
// 		// todo if no progress in io copy then cancel the request
// 		n, err := io.CopyN(commonWriter, file, size)
// 		if err != nil {
// 			return err
// 		}
// 		logrus.Infof("io.copy %v bytes", n)
// 		return nil
// 	})
// 	if err := g.Wait(); err != nil {
// 		logrus.Errorf("io.copy failed: %v", err)
// 		c.Status(http.StatusInternalServerError)
// 		return
// 	}
// 	// wait for minioUpload
// 	select {
// 	case <-minioUploadDone:
// 		logrus.Info("file uploaded to minio")
// 	case err := <-minioUploadErrs:
// 		logrus.Errorf("file upload failed: %v", err)
// 		c.Status(http.StatusInternalServerError)
// 		return
// 	}

// 	// wait for thumbnail
// 	select {
// 	case <-thumbnailGeneratorDone:
// 		logrus.Info("thumbnail generated")
// 	case err := <-thumbnailGeneratorErrs:
// 		logrus.Errorf("file upload failed: %v", err)
// 		c.Status(http.StatusInternalServerError)
// 		return
// 	}
// 	c.Status(http.StatusOK)
// 	return
// }
