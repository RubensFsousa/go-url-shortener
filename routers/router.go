package routers

import (
	"github.com/RubensFsousa/go-url-shortener/handler"
	"github.com/gin-gonic/gin"
)

func InitRouters() {
	router := gin.Default()

	v1 := router.Group("/api/url")
	{
		v1.POST("/", handler.CoderUrlHandler)
		v1.GET("/:shortUrl", handler.DecoderUrlHandler)
	}

	router.Run()
}
