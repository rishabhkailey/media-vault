package api

import (
	"fmt"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	v1Api "github.com/rishabhkailey/media-vault/internal/api/v1"
	"github.com/rishabhkailey/media-vault/internal/api/website"
	"github.com/rishabhkailey/media-vault/internal/config"
	"github.com/sirupsen/logrus"
)

func NewRouter(v1ApiServer *v1Api.Server, websiteHandler *website.WebsiteHandler, config config.Config) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(v1Api.ErrorHandler)

	router.NoRoute(websiteHandler.ServeWebsite)

	v1 := router.Group("/v1")
	{
		// refresh session has its own auth validation
		v1.POST("/refreshSession", v1ApiServer.RefreshSession)

		// session based and fallback to bearer token
		v1UserProtected := v1.Group("/")
		v1UserProtected.Use(v1ApiServer.UserAuthMiddleware)
		{
			v1UserProtected.POST("/upload", v1ApiServer.InitChunkUpload)
			v1UserProtected.POST("/upload/:upload_request_id/chunk", v1ApiServer.UploadChunk)
			v1UserProtected.POST("/upload/:upload_request_id/finish", v1ApiServer.FinishChunkUpload)
			v1UserProtected.POST("/upload/:upload_request_id/thumbnail", v1ApiServer.UploadThumbnail)
			v1UserProtected.GET("/media-list", v1ApiServer.MediaList)
			v1UserProtected.GET("/media/:media_id", v1ApiServer.GetMedia)
			v1UserProtected.GET("/search", v1ApiServer.Search)
			v1UserProtected.DELETE("/media/:media_id", v1ApiServer.DeleteSingleMedia)
			v1UserProtected.DELETE("/media", v1ApiServer.DeleteMedia)
			v1UserProtected.POST("/album", v1ApiServer.CreateAlbum)
			v1UserProtected.GET("/albums", v1ApiServer.GetAlbums)
			v1UserProtected.GET("/album/:album_id", v1ApiServer.GetAlbum)
			v1UserProtected.PATCH("/album/:album_id", v1ApiServer.PatchAlbum)
			v1UserProtected.DELETE("/album/:album_id", v1ApiServer.DeleteAlbum)
			v1UserProtected.GET("/album/:album_id/media", v1ApiServer.GetAlubmMedia)
			v1UserProtected.POST("/album/:album_id/media", v1ApiServer.AlbumAddMedia)
			v1UserProtected.DELETE("/album/:album_id/media", v1ApiServer.RemoveAlbumMedia)
			v1UserProtected.GET("/user-info", v1ApiServer.GetUserInfo)
			v1UserProtected.POST("/user-info", v1ApiServer.PostUserInfo)
		}

		// session based only
		v1FileAccessProtected := v1.Group("/")
		v1FileAccessProtected.Use(v1ApiServer.SessionBasedMediaFileAuthMiddleware)
		{
			v1FileAccessProtected.GET("/file/:file_name", v1ApiServer.GetMediaFile)
			v1FileAccessProtected.GET("/file/:file_name/thumbnail", v1ApiServer.GetThumbnailFile)
		}

		// bearer token only
		v1ProtectedBearerOnly := v1.Group("/")
		v1ProtectedBearerOnly.Use(v1ApiServer.UserAuthMiddleware)
		{
			v1ProtectedBearerOnly.POST("/terminate-session", v1ApiServer.TerminateSession)
		}
	}
	public := v1.Group("/public")
	{
		public.GET("/spa-config", v1ApiServer.GetSpaConfig)
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
	websiteHandler := website.NewWebsiteHandler(config.WebUIConfig.Directory)

	router, err := NewRouter(v1ApiServer, websiteHandler, *config)
	if err != nil {
		return fmt.Errorf("[Server.Start] failed to create router: %w", err)
	}
	return router.Run(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
}
