package v1

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/awnumar/memguard"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/sirupsen/logrus"
)

const tempEncryptedDirectory = "./tempEncryptedStorage"

func (server *Server) TestEncryptedFileSave(c *gin.Context) {

	memguard.CatchInterrupt()
	defer memguard.Purge()
	secureKey := memguard.NewEnclaveRandom(16)
	keyBuf, err := secureKey.Open()
	if err != nil {
		logrus.Errorf("securekey.Open failed %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer keyBuf.Destroy()

	bucketname := "test"
	err = utils.CreateBucketIfMissing(c.Request.Context(), *server.Minio, bucketname)
	if err != nil {
		logrus.WithField(
			"error", err,
		).Error("bucket creation failed")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// objectName := "my-objectname"
	requestId, err := uuid.NewRandom()
	if err != nil {
		logrus.Errorf("failed to generate random uuid: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	tempDirectory := path.Join(tempEncryptedDirectory, requestId.String())
	if os.MkdirAll(tempDirectory, 0750) != nil {
		logrus.Errorf("failed to create temp directory: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// defer os.RemoveAll(tempDirectory)
	reqReader, err := c.Request.MultipartReader()
	if err != nil {
		logrus.Errorf("failed to get request reader: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var name string
	var size int64
	var file *os.File
	for {
		part, err := reqReader.NextPart()
		if err == io.EOF {
			break
		}
		switch n := part.FormName(); n {
		case "name":
			{
				b, err := io.ReadAll(part)
				if err != nil && !errors.Is(err, io.EOF) {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				name = string(b)
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
				filePath := path.Join(tempDirectory, "encrypted-file")
				var err error
				file, err = os.Create(filePath)
				if err != nil {
					logrus.Errorf("%v file creation failed: %v", filePath, err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				defer file.Close()
				var n int64
				n, err = copyToEncryptedFile(file, part, keyBuf.Bytes())
				if err != nil {
					logrus.Errorf("%v file creation failed: %v", filePath, err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				logrus.Infof("succesfully copied %v bytes", n)
			}
		}
	}
	if file == nil || size == 0 || len(name) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// file upload
	{
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			logrus.Errorf("encrypted file seek failed: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		uploadedInfo, err := server.Minio.PutObject(c.Request.Context(), bucketname, name, createEncryptedReader(file, keyBuf.Bytes()), size, minio.PutObjectOptions{})
		if err != nil {
			logrus.Errorf("%v file upload failed: %v", name, err)
		}
		logrus.WithField("info", uploadedInfo).Info("file uploaded")

	}

	// thumbnail generate and upload
	{
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			logrus.Errorf("encrypted file seek failed: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		thumbnailName := fmt.Sprintf(".thumb-%s.%s", name, "png")
		{
			ns := strings.Split(name, ".")
			if len(name) > 1 {
				thumbnailName = fmt.Sprintf(".thumb-%s.%s", ns[0], "png")
			}
		}

		thumbnailBytes, err := utils.GenerateThumbnail(createEncryptedReader(file, keyBuf.Bytes()), 300)
		if err != nil {
			logrus.Errorf("generate thumbnail failed: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		bytes.NewReader(thumbnailBytes)
		uploadedInfo, err := server.Minio.PutObject(c.Request.Context(), bucketname, thumbnailName, bytes.NewReader(thumbnailBytes), int64(len(thumbnailBytes)), minio.PutObjectOptions{})
		if err != nil {
			logrus.Errorf("%v file upload failed: %v", thumbnailName, err)
		}
		logrus.WithField("info", uploadedInfo).Info("file uploaded")
	}

	c.Status(http.StatusOK)
}

func copyToEncryptedFile(file io.Writer, r io.Reader, encryptionKey []byte) (n int64, err error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return
	}
	var iv [aes.BlockSize]byte

	encryptedWriter := &cipher.StreamWriter{S: cipher.NewOFB(block, iv[:]), W: file}
	if n, err = io.Copy(encryptedWriter, r); err != nil {
		err = fmt.Errorf("encrypt write failed: %w", err)
		return
	}
	return
}

func createEncryptedReader(r io.Reader, encryptionKey []byte) (encryptedReader io.Reader) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return
	}
	var iv [aes.BlockSize]byte
	return cipher.StreamReader{S: cipher.NewOFB(block, iv[:]), R: r}
}
