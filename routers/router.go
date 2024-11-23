package routers

import (
	"github.com/RubensFsousa/go-url-shortener/handler"
	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	v1 := r.Group("/api/url")
	{
		v1.POST("/", handler.CoderUrlHandler)
		v1.GET("/:shortUrl", handler.DecoderUrlHandler)
	}
}
