package internal

import (
	"go.uber.org/zap"
	"golang-web-template/internal/config"
	"log"
	"strings"
)

// AppContext holds all the necessary dependencies that are required to
// start app.Runnable.
type AppContext struct {
	Logger *zap.Logger
	//DbRepository DbRepository
}

func InitAppContext(cfg *config.AppConfig) *AppContext {
	logger := initLogger(cfg.Env)

	return &AppContext{
		Logger: logger,
	}
}

func initLogger(env string) *zap.Logger {
	var logger *zap.Logger
	var err error
	if strings.ToLower(env) == "prod" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("failed to create zap logger with error: %v", err)
	}
	return logger
}
