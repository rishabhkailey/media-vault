package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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
	thumbnailName := fmt.Sprintf(".thumb-%s.%s", fileName, "png")
	{
		name := strings.Split(fileName, ".")
		if len(name) > 1 {
			thumbnailName = fmt.Sprintf(".thumb-%s.%s", name[0], "png")
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
	thumbnailBytes, err := utils.GenerateThumbnail(file, 300)
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

// this will not work
// nextPart closes the current part so this will not work
func (server *Server) TestStreamVideoUploadWithThumbnail(c *gin.Context) {

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
	reqReader, err := c.Request.MultipartReader()
	if err != nil {
		logrus.Errorf("failed to get request reader: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var fileName string
	var size int64
	var file *multipart.Part
	for {
		// this closes the current part so this will not work
		part, err := reqReader.NextPart()
		if err == io.EOF {
			break
		}
		// defer part.Close()
		switch n := part.FormName(); n {
		case "name":
			{
				b, err := io.ReadAll(part)
				if err != nil && !errors.Is(err, io.EOF) {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				fileName = string(b)
			}
		case "size":
			{
				b, err := io.ReadAll(part)
				if err != nil && !errors.Is(err, io.EOF) {
					// todo bad request?
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				size, err = strconv.ParseInt(string(b), 10, 64)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
			}
		case "file":
			{
				file = part
			}
		}
	}
	if file == nil || size == 0 || len(fileName) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	logrus.WithFields(logrus.Fields{
		"file": file,
	}).Info("file request")

	thumbnailName := fmt.Sprintf(".thumb-%s.%s", fileName, "png")
	{
		name := strings.Split(fileName, ".")
		if len(name) > 1 {
			thumbnailName = fmt.Sprintf(".thumb-%s.%s", name[0], "png")
		}
	}

	// ffmpegReader, ffmpegWriter := io.Pipe()
	// minioReader := io.TeeReader(file, ffmpegWriter)

	// todo check error groups instead of using channel if that makes sense. as we don't have any data in done channel
	minioUploadErrs := make(chan error)
	thumbnailGeneratorErrs := make(chan error)
	minioUploadDone := make(chan bool)
	thumbnailGeneratorDone := make(chan bool)

	go tempFileSave(file, "ffmpeg.tmp")
	// go tempFileSave(minioReader, "minio.tmp")
	_ = thumbnailName
	// go server.uploadFileToMinio(c.Request.Context(), bucketname, fileName, minioReader, int64(size), minioUploadDone, minioUploadErrs)
	// go server.generateAndUploadThumbnail(c.Request.Context(), bucketname, thumbnailName, ffmpegReader, size, thumbnailGeneratorDone, thumbnailGeneratorErrs)

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

func tempFileSave(r io.Reader, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	logrus.Infof("copied %v bytes in %v", n, file)
	return nil
}
