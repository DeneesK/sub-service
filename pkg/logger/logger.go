package logger

import (
	"go.uber.org/zap"
)

const (
	debug = "debug"
	info  = "info"
)

func NewLogger(logLevel string) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	switch logLevel {
	case debug:
		logger, err = zap.NewDevelopment()
	case info:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		logger.Fatal("failed to initialized new logger", zap.String("err", err.Error()))
	}

	sugar := logger.Sugar()
	return sugar
}
