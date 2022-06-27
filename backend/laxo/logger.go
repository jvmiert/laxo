package laxo

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	Zap *zap.Logger
	*zap.SugaredLogger
}

func NewLogger() *Logger {
	config := zap.NewDevelopmentConfig()

	zapLogger, err := config.Build()
	if err != nil {
		log.Fatal("Could not build logger")
	}

	sugar := zapLogger.Sugar()

	logger := &Logger{
		zapLogger,
		sugar,
	}

	return logger
}
