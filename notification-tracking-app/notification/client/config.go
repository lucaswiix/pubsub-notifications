package client

import "go.uber.org/zap/zapcore"

type RabbitMQConfig struct {
	Address   string       `mapstructure:"rabbitmq_addr"`
	QueueName string       `mapstructure:"queue_name"`
	Logger    LoggerConfig `mapstructure:"logger"`
}

type LoggerConfig struct {
	Level zapcore.Level
}

var DefaultRabbitConfig = RabbitMQConfig{
	Logger: LoggerConfig{
		Level: zapcore.InfoLevel,
	},
}
