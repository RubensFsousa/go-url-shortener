package main

import (
	"github.com/RubensFsousa/go-url-shortener/config"
	_ "github.com/RubensFsousa/go-url-shortener/docs"
	"github.com/RubensFsousa/go-url-shortener/routers"
)

var (
	logger *config.Logger
)

// @title           Shortener url go api
// @version         1.0
// @description     api to store and shorten urls
// @BasePath        /api
func main() {
	logger = config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Errorf("Initialization error: %v", err)
		return
	}

	routers.InitRouters()
}
