package repository

import (
	"meli/notifications/utils"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	DB *redis.Client
	CH *amqp.Channel
)

func InitDatabase() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("input.redis.addr"),
		Password: viper.GetString("input.redis.password"),
		DB:       0,
	})

	DB = rdb
	return nil
}

func NewRabbitMQ() error {
	topicUrl := viper.GetString("input.rabbitmq.topic_url")
	topicName := viper.GetString("input.rabbitmq.topic")
	conn, err := amqp.Dial(topicUrl)
	if err != nil {
		utils.Log.Error("Failed to connect to rabbitmq", zap.Error(err))
		return err
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		utils.Log.Error("Failed to open a channel", zap.Error(err))
		return err
	}
	// defer ch.Close()

	headers := make(amqp.Table)
	headers["x-delayed-type"] = "direct"
	err = ch.ExchangeDeclare("notifications", "x-delayed-message", false, false, false, true, headers)
	if err != nil {
		utils.Log.Error("Failed to declare a exchange", zap.Error(err))
		return err
	}

	_, err = ch.QueueDeclare(
		topicName, // name
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		false,
		nil, // arguments
	)
	if err != nil {
		utils.Log.Error("Failed to declare a queue", zap.Error(err))
		return err
	}

	err = ch.QueueBind(
		topicName,
		"",
		"notifications",
		false,
		nil)

	if err != nil {
		utils.Log.Error("Failed to bind to the queue", zap.Error(err))
		return err
	}

	CH = ch
	return nil
}
