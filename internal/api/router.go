package api

import (
	"fmt"
	"net/http/pprof"
	"path"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	v1Api "github.com/rishabhkailey/media-service/internal/api/v1"
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/sirupsen/logrus"
)

func NewRouter(v1ApiServer *v1Api.Server, config config.Config) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(v1Api.ErrorHandler)
	// store := cookie.NewStore([]byte(config.Session.Secret))
	// // session.SetCookieName("media_service")
	// router.Use(sessions.Sessions("media_service", store))

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
		v1UserProtected := v1.Group("/")
		v1UserProtected.Use(v1ApiServer.UserAuthMiddleware)
		{
			// todo gitlab style api endpoints?
			// POST /v1/upload (init)
			// POST /v1/upload/:upload_request_id/chunk (upload chunk)
			// POST /v1/upload/:upload_request_id/thumbnail (upload thumbnail)
			// POST /v1/upload/:upload_request_id/finish (finish upload)
			v1UserProtected.POST("/initChunkUpload", v1ApiServer.InitChunkUpload)
			v1UserProtected.POST("/uploadChunk", v1ApiServer.UploadChunk)
			v1UserProtected.POST("/finishChunkUpload", v1ApiServer.FinishChunkUpload)
			v1UserProtected.POST("/uploadThumbnail", v1ApiServer.UploadThumbnail)
			v1UserProtected.GET("/mediaList", v1ApiServer.MediaList)
			v1UserProtected.GET("/media/:media_id", v1ApiServer.GetMedia)
			v1UserProtected.GET("/search", v1ApiServer.Search)
			v1UserProtected.DELETE("/media/:media_id", v1ApiServer.DeleteMedia)
			v1UserProtected.POST("/album", v1ApiServer.CreateAlbum)
			v1UserProtected.GET("/albums", v1ApiServer.GetAlbums)
			v1UserProtected.GET("/album/:albumID", v1ApiServer.GetAlbum)
			v1UserProtected.PATCH("/album/:albumID", v1ApiServer.PatchAlbum)
			v1UserProtected.DELETE("/album/:albumID", v1ApiServer.DeleteAlbum)
			v1UserProtected.GET("/album/:albumID/media", v1ApiServer.GetAlubmMedia)
			v1UserProtected.POST("/album/:albumID/media", v1ApiServer.AlbumAddMedia)
			v1UserProtected.DELETE("/album/:albumID/media", v1ApiServer.RemoveAlbumMedia)
			v1UserProtected.GET("/user-info", v1ApiServer.GetUserInfo)
			v1UserProtected.POST("/user-info", v1ApiServer.PostUserInfo)
		}

		// session based only
		v1FileAccessProtected := v1.Group("/")
		v1FileAccessProtected.Use(v1ApiServer.SessionBasedMediaFileAuthMiddleware)
		{
			v1FileAccessProtected.GET("/file/:fileName", v1ApiServer.GetMediaFile)
			v1FileAccessProtected.GET("/file/:fileName/thumbnail", v1ApiServer.GetThumbnailFile)
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
