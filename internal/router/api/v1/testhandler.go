package v1

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (server *Server) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hi",
		"success": true,
	})
}

// no ecoding, no encryption, no parts
func (server *Server) TestGetVideo(c *gin.Context) {
	object, err := server.Minio.GetObject(c.Request.Context(), "test", "test.mp4", minio.GetObjectOptions{})
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	objInfo, err := object.Stat()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))
	c.Header("Content-Type", "video/mp4")
	c.Header("Connection", "keep-alive")
	// sending whole video at once
	c.Header("Content-Range", fmt.Sprintf("bytes 0-%d/%d", objInfo.Size-1, objInfo.Size))
	c.Header("Accept-Ranges", "bytes")
	// c.SSEvent()
	c.Status(http.StatusOK)
	c.Stream(func(w io.Writer) bool {
		n, err := io.Copy(w, object)
		logrus.WithField("bytes", n).Info("sent")
		if err != nil {
			logrus.Error(err)
		}
		return false
	})
}

func (server *Server) TestDownload(c *gin.Context) {
	fileName := c.Request.FormValue("file")
	if len(fileName) == 0 {
		logrus.WithField(
			"file", fileName,
		).Error("not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// object, err := server.Minio.GetObject(c.Request.Context(), "test", fileName, minio.GetObjectOptions{})
	object, err := server.MinioObjectCache.Get(c.Request.Context(), "test", fileName)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	objectInfo, err := object.Stat()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Length", strconv.FormatInt(objectInfo.Size, 10))
	c.Header("Connection", "keep-alive")
	logrus.Info(io.Copy(c.Writer, object))
}

// todo cache headers cache-control, x-cache ...
// todo io.CopyBuffer

// todo dumb cache
// var cache map[context.Context]*minio.Object = make(map[context.Context]*minio.Object)

// func (server *Server) getObject(c context.Context, bucketName, objectName string) (*minio.Object, error) {
// 	var object minio.Object
// 	if obj, ok := cache[c]; ok {
// 	}

// 	return nil, nil
// }

// range
// https://stackoverflow.com/questions/3303029/http-range-header
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
// https://stackoverflow.com/questions/13043816/html5-video-and-partial-range-http-requests
func (server *Server) TestGetVideoWithRange(c *gin.Context) {
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
	// object, err := server.Minio.GetObject(c.Request.Context(), "test", fileName, minio.GetObjectOptions{})
	object, err := server.MinioObjectCache.Get(c.Request.Context(), "test", fileName)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// todo support for multiple ranges
	if parsedRangeHeader == nil || len(parsedRangeHeader.ranges) != 1 {
		server.TestGetVideoFirstRequest(c, object)
		// server.TestGetRangeVideo(c, Range{
		// 	start: 0,
		// 	end:   -1,
		// })
		return
	}
	server.TestGetRangeVideo(c, parsedRangeHeader.ranges[0], object)
}

func (server *Server) TestGetVideoFirstRequest(c *gin.Context, object *minio.Object) {
	objInfo, err := object.Stat()
	if err != nil {
		logrus.Error(err)
	}
	// this is for giving client the hint that response is a video file
	contentLength := objInfo.Size
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Content-Type", "video/mp4")
	c.Header("Connection", "keep-alive")
	c.Header("Accept-Ranges", "bytes")
	c.Status(http.StatusPartialContent)
}

// todo browsers which don't support range requests
// todo what to do on first request without range
// https://vjs.zencdn.net/v/oceans.mp4 this return a 200 response with content length only?
// if range end not provided
// const defaultRangeSize int64 = 1000000 // 1mb
func (server *Server) TestGetRangeVideo(c *gin.Context, r Range, object *minio.Object) {

	objInfo, err := object.Stat()
	if err != nil {
		logrus.Error(err)
	}
	if objInfo.Size == 0 {
		// todo gracefully handle to this for empty files
		c.Status(http.StatusInternalServerError)
		return
	}
	if r.end == -1 {
		r.end = r.start + defaultRangeSize
	}
	if r.end > objInfo.Size-1 {
		r.end = objInfo.Size - 1
	}
	contentLength := r.end - r.start + 1
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Content-Type", "video/mp4")
	c.Header("Connection", "keep-alive")
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", r.start, r.end, objInfo.Size))
	c.Header("Accept-Ranges", "bytes")
	// c.SSEvent()
	// todo use of stream?
	logrus.WithField("range", r).Info("request received")
	_, err = object.Seek(r.start, 0)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusPartialContent)
	n, err := io.CopyN(c.Writer, object, contentLength)
	logrus.WithField("bytes", n).Info("sent")
	if err != nil {
		// todo this will not helm i guess, status code set earlier will be sent when we start copying the data
		c.Status(http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	return
}

func (server *Server) TestGetFullVideo(c *gin.Context) {
	object, err := server.Minio.GetObject(c.Request.Context(), "test", "test.mp4", minio.GetObjectOptions{})
	if err != nil {
		logrus.Error(err)
	}
	objInfo, err := object.Stat()
	if err != nil {
		logrus.Error(err)
	}

	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))
	c.Header("Content-Type", "video/mp4")
	c.Header("Connection", "keep-alive")
	// sending whole video at once
	c.Header("Content-Range", fmt.Sprintf("bytes 0-%d/%d", objInfo.Size-1, objInfo.Size))
	c.Header("Accept-Ranges", "bytes")
	// c.SSEvent()
	c.Status(http.StatusPartialContent)
	c.Stream(func(w io.Writer) bool {
		n, err := io.Copy(w, object)
		logrus.WithField("bytes", n).Info("sent")
		if err != nil {
			logrus.Error(err)
		}
		return false
	})
}

// range is inclusive
type Range struct {
	start int64 // start will always be provided
	end   int64 // end = -1 if not provided
}
type RangeHeader struct {
	unit   string
	ranges []Range
}

func parseRangeHeader(value string) (*RangeHeader, error) {
	s := strings.Split(value, "=")
	if len(s) != 2 {
		return nil, fmt.Errorf("invalid range header")
	}
	unit := s[0]
	var ranges []Range
	rangesStr := strings.Split(s[1], ",")
	for _, r := range rangesStr {
		r, err := parseRange(r)
		if err != nil {
			return nil, fmt.Errorf("invalid range header: %w", err)
		}
		ranges = append(ranges, *r)
	}
	return &RangeHeader{
		unit:   unit,
		ranges: ranges,
	}, nil
}

func parseRange(r string) (*Range, error) {
	r = strings.TrimSpace(r)
	separatorIndex := strings.Index(r, "-")
	// validations
	if separatorIndex == -1 || separatorIndex == 0 {
		return nil, fmt.Errorf("invalid range: %v", r)
	}
	rangeStartEndArr := strings.Split(r, "-")
	startStr := rangeStartEndArr[0]
	endStr := rangeStartEndArr[1]
	if len(rangeStartEndArr) > 2 {
		return nil, fmt.Errorf("invalid range: %v", r)
	}
	// range start
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid range: %w", err)
	}

	// range end
	end, err := strconv.ParseInt(rangeStartEndArr[1], 10, 64)
	if err != nil {
		if len(endStr) != 0 {
			return nil, fmt.Errorf("invalid range: %w", err)
		}
		end = -1
	}
	return &Range{
		start: int64(start),
		end:   int64(end),
	}, nil
}
