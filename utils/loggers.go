package utils

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {

	config := zap.NewProductionConfig()

	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		log.Fatalln(err)
	}

	Log = logger

}
