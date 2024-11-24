package main

import (
	"github.com/RubensFsousa/go-url-shortener/config"
	"github.com/RubensFsousa/go-url-shortener/routers"
)

var (
	logger *config.Logger
)

func main() {
	logger := config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Errorf("Initialization error: %v", err)
		return
	}

	routers.InitRouters()
}
