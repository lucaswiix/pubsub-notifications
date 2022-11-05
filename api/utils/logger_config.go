package utils

import "go.uber.org/zap/zapcore"

type Config struct {
	Logger LoggerConfig `mapstructure:"logger"`
}

type LoggerConfig struct {
	Level zapcore.Level
}
