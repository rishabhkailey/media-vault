package router

import (
	"fmt"
	"net/http/pprof"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rishabhkailey/media-service/internal/config"
	v1Api "github.com/rishabhkailey/media-service/internal/router/api/v1"
	"github.com/sirupsen/logrus"
)

func NewRouter(v1ApiServer *v1Api.Server, config config.Config) (*gin.Engine, error) {
	router := gin.Default()
	store := cookie.NewStore([]byte(v1ApiServer.Config.Session.Secret))
	router.Use(sessions.Sessions("mysession", store))

	// todo - check if there is any security risk of doing this
	router.Use(
		static.Serve("/assets", static.LocalFile(path.Join(config.WebUIConfig.Directory, "assets"), false)),
		static.Serve("/login", static.LocalFile(config.WebUIConfig.Directory, false)),
		static.Serve("/signup", static.LocalFile(config.WebUIConfig.Directory, false)),
		static.Serve("/consentscreen", static.LocalFile(config.WebUIConfig.Directory, false)),
	)

	v1 := router.Group("/v1")
	{
		v1.GET("/login", v1ApiServer.LoginHandler)
		v1.POST("/logout", v1ApiServer.LogoutHandler)
		v1.GET("/authorize", v1ApiServer.AuthHandler)
		v1.GET("/testGetVideo", v1ApiServer.TestGetVideo)
		v1.GET("/testGetVideoWithRange", v1ApiServer.TestGetVideoWithRange)
		v1.GET("/testDownloadFileWithRange", v1ApiServer.TestGetVideoWithRange)
		v1.GET("/testDownload", v1ApiServer.TestDownload)
		v1.GET("/testGetEncryptedVideoWithRange", v1ApiServer.TestGetVideoWithRange)
		v1.GET("/testGetVideoWithRange/test.mp4", v1ApiServer.TestGetVideoWithRange)
		v1.POST("/testNormalUpload", v1ApiServer.TestNormalUpload)
		v1.POST("/testEncryptedUpload", v1ApiServer.TestEncryptedUpload)
		v1.GET("/testGetEncryptedVideo", v1ApiServer.TestGetEncryptedVideo)
		v1.GET("/testGetEncryptedImage", v1ApiServer.TestGetEncryptedImage)
		v1.POST("/testVideoUploadWithThumbnail", v1ApiServer.TestVideoUploadWithThumbnail)
		v1.POST("/testEncryptedFileSave", v1ApiServer.TestEncryptedFileSave)
		v1.POST("/testStreamVideoUploadWithThumbnail", v1ApiServer.TestStreamVideoUploadWithThumbnail)
		v1.POST("/initChunkUpload", v1ApiServer.InitChunkUpload)
		v1.POST("/uploadChunk", v1ApiServer.UploadChunk)
		v1.POST("/finishChunkUpload", v1ApiServer.FinishChunkUpload)
	}
	// authorized endpoints
	authorized := router.Group("/")
	// is it good to have redirects at backend?
	// we will have the AuthMiddleware just for authentication check and will move the redirect part to /login method
	// add current /login will be renamed to afterLogin or postlogin
	authorized.Use(v1ApiServer.AuthMiddleware)
	{
		v1Authorized := authorized.Group("/v1")
		{
			v1Authorized.GET("/userinfo", v1ApiServer.UserInfo)
			v1Authorized.GET("/testProtectedGetEncryptedImage", v1ApiServer.TestGetEncryptedImage)
		}
	}

	debug := router.Group("/debug")
	{
		debug.GET("/pprof/:profile", func(c *gin.Context) {
			profile := c.Param("profile")
			if profile == "profile" {
				pprof.Profile(c.Writer, c.Request)
				return
			}
			handler := pprof.Handler(profile)
			handler.ServeHTTP(c.Writer, c.Request)
		})
	}
	return router, nil
}

func Start() error {
	config, err := config.GetConfig()
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	logrus.Info(config)

	v1ApiServer, err := v1Api.NewServer(config)
	if err != nil {
		return fmt.Errorf("[Server.Start] failed to create Server: %w", err)
	}
	router, err := NewRouter(v1ApiServer, *config)
	if err != nil {
		return fmt.Errorf("[Server.Start] failed to create router: %w", err)
	}
	return router.Run(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
}
