package utils

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	logger = logger.With(zap.String("app", "pubsub-notification-api")).With(zap.String("environment", "psm"))

	Log = logger

}
