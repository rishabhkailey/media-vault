package website

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WebsiteHandler struct {
	fileSystem http.FileSystem
}

func NewWebsiteHandler(directory string) *WebsiteHandler {
	return &WebsiteHandler{
		fileSystem: http.FS(os.DirFS(directory)),
	}
}

func (wh *WebsiteHandler) ServeWebsite(c *gin.Context) {
	filePath := c.Request.URL.Path
	file, err := wh.fileSystem.Open(filePath)
	if errors.Is(err, fs.ErrNotExist) || filePath == "/" {
		file, err = wh.fileSystem.Open("index.html")
	}
	if err != nil {
		c.Error(fmt.Errorf("[static.Serve] unable to open %s file: %w", filePath, err))
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.Error(fmt.Errorf("[static.Serve] unable to get %s file info: %w", filePath, err))
		return
	}
	logrus.Infof("[ServeWebsite]: serving %s", fileInfo.Name())
	c.Header("Service-Worker-Allowed", "/")
	http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)
}
