package main

import (
	"github.com/DeneesK/sub-service/internal/app"
	"github.com/DeneesK/sub-service/internal/config"
	"github.com/DeneesK/sub-service/internal/db"
	"github.com/DeneesK/sub-service/internal/service"
	"github.com/DeneesK/sub-service/pkg/logger"
)

func main() {
	conf := config.MustLoad()

	log := logger.NewLogger(conf.LogLevel)

	db, err := db.InitDBConnection(
		conf.MigrationPath,
		conf.DBHost, conf.DBPort,
		conf.DBUser, conf.DBPassword,
		conf.DBName, conf.DBSSLMode,
	)

	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	log.Info("DB initialized successfully")
	subService := service.NewSubscriptionService(db, log)

	a := app.NewApp(conf.ServerAddr, log, subService)
	a.Run()
}
