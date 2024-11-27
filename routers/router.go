package routers

import (
	"os"

	"github.com/RubensFsousa/go-url-shortener/handler"
	"github.com/gin-gonic/gin"
)

func InitRouters() {
	router := gin.Default()
	handler.InitializeHandler()

	v1 := router.Group("/api/url")
	{
		v1.POST("/", handler.CoderUrlHandler)
		v1.GET("/:codedUrl", handler.DecoderUrlHandler)
	}

	router.Run(os.Getenv("PORT"))
}
