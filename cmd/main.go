package main

import (
	"github.com/DeneesK/sub-service/internal/config"
	"github.com/DeneesK/sub-service/internal/db"
	"github.com/DeneesK/sub-service/pkg/logger"
)

func main() {
	conf := config.MustLoad()

	db := db.NewDBConnection(conf.MigrationPath,
		conf.DBHost, conf.DBPort,
		conf.DBUser, conf.DBPassword,
		conf.DBName, conf.DBSSLMode)

	log := logger.NewLogger(conf.LogLevel)

}
