package web

import (
	"imager/internal/service"
	"imager/internal/web/handlers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
)

func SetRoutes(router *ginext.Engine, templates string, srv *service.Service) {
	router.Delims("__", "__").LoadHTMLGlob(templates)

	router.Use(ginext.Logger())

	router.GET("/", handlers.MainPage)
	router.GET("//upload-file", handlers.UploadPage)

	router.POST("/upload", handlers.UploadFile(srv))
	router.GET("/image/:id", handlers.GetFile(srv))
	router.DELETE("/image/:id", handlers.DeleteFile(srv))

	router.StaticFS("/processed", gin.Dir("processed_pictures/", true))
}

func NewServer(router *ginext.Engine, cfg *config.Config) *http.Server {
	rt, err := time.ParseDuration(cfg.GetString("read_timeout"))
	if err != nil {
		panic(err)
	}
	wt, err := time.ParseDuration(cfg.GetString("write_timeout"))
	if err != nil {
		panic(err)
	}

	return &http.Server{
		Handler:      router,
		ReadTimeout:  rt,
		WriteTimeout: wt,
	}
}
