package main

import (
	"github.com/DeneesK/sub-service/internal/config"
	"github.com/DeneesK/sub-service/pkg/logger"
)

func main() {
	conf := config.MustLoad()

	log := logger.NewLogger(conf.LogLevel)

}
