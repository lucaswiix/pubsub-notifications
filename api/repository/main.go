package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/lucaswiix/meli/notifications/utils"

	"github.com/lucaswiix/meli/notifications/pkg/orm"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	DB                 dynamodbiface.DynamoDBAPI
	CH                 *amqp.Channel
	UsersTable         = "users"
	NotificationsTable = "notifications"
)

func InitDatabase() error {

	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(viper.GetString("input.dynamodb.addr")),
		Region:   aws.String("us-east-1"),
	},
	)
	if err != nil {
		utils.Log.Error("error to connect to aws session", zap.Error(err))
	}

	svc := dynamodb.New(sess)

	config := orm.NewDynamoDBORM(svc, UsersTable)

	config.CreateTable(UsersTable)

	DB = dynamodbiface.DynamoDBAPI(svc)
	return nil
}

func NewRabbitMQ() error {
	topicUrl := viper.GetString("input.rabbitmq.topic_url")
	topicName := viper.GetString("input.rabbitmq.topic")
	utils.Log.Info(topicUrl)
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
		true,
		false,
		false,
		false,
		nil,
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
