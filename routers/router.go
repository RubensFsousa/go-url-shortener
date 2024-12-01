package routers

import (
	"os"

	"github.com/RubensFsousa/go-url-shortener/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouters() {
	router := gin.Default()
	handler.InitializeHandler()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/url")
	{
		v1.POST("/", handler.CoderUrlHandler)
		v1.GET("/:codedUrl", handler.DecoderUrlHandler)
	}

	port := ":" + os.Getenv("PORT")

	router.Run(port)
}
