package main

import (
	"github.com/RubensFsousa/go-url-shortener/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.InitRouters(r)
	r.Run()
}
