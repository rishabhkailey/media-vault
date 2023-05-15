package api

import (
	"fmt"
	"net/http/pprof"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	v1Api "github.com/rishabhkailey/media-service/internal/api/v1"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/sirupsen/logrus"
)

func NewRouter(v1ApiServer *v1Api.Server, config config.Config) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(v1Api.ErrorHandler)
	store := cookie.NewStore([]byte(config.Session.Secret))
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
		// refresh session has its own auth validation
		v1.POST("/refreshSession", v1ApiServer.RefreshSession)

		// session based and fallback to bearer token
		userProtected := v1.Group("/")
		userProtected.Use(v1ApiServer.UserAuthMiddleware)
		{
			v1UserProtected := userProtected.Group("/")
			{
				v1UserProtected.POST("/initChunkUpload", v1ApiServer.InitChunkUpload)
				v1UserProtected.POST("/uploadChunk", v1ApiServer.UploadChunk)
				v1UserProtected.POST("/finishChunkUpload", v1ApiServer.FinishChunkUpload)
				v1UserProtected.POST("/uploadThumbnail", v1ApiServer.UploadThumbnail)
				v1UserProtected.GET("/mediaList", v1ApiServer.MediaList)
				v1UserProtected.GET("/search", v1ApiServer.Search)
			}
		}

		// session based only
		v1FileAccessProtected := v1.Group("/")
		v1FileAccessProtected.Use(v1ApiServer.SessionBasedMediaAuthMiddleware)
		{
			v1FileAccessProtected.GET("/media/:fileName", v1ApiServer.GetMedia)
			v1FileAccessProtected.GET("/thumbnail/:fileName", v1ApiServer.GetThumbnail)
		}

		// bearer token only
		v1ProtectedBearerOnly := v1.Group("/")
		v1ProtectedBearerOnly.Use(v1ApiServer.UserAuthMiddleware)
		{
			v1ProtectedBearerOnly.POST("/terminateSession", v1ApiServer.TerminateSession)
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
