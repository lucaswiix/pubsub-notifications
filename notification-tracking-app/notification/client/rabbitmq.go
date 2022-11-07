package client

import (
	"context"
	"errors"

	"github.com/lucaswiix/notifications-tracking-app/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type rabbitmqClient struct {
	conn               *amqp.Connection
	ch                 *amqp.Channel
	connString         string
	notificationStatus <-chan amqp.Delivery
}

func NewRabbitMQClient(connectionString, queueName string) (*rabbitmqClient, error) {
	c := &rabbitmqClient{}
	var err error

	c.conn, err = amqp.Dial(connectionString)
	if err != nil {
		utils.Log.Error("error on try to connect with amqp", zap.Error(err))
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		utils.Log.Error("error on try to make a amqp channel", zap.Error(err))
		return nil, err
	}

	err = c.configureQueue(queueName)

	return c, err
}

func (c *rabbitmqClient) ConsumeByUserID(ctx context.Context, userID, nType string) ([]byte, error) {
	for msg := range c.notificationStatus {
		if msg.MessageId == userID && msg.Type == nType {
			utils.Log.Debug("receive notification", zap.String("to user id", msg.MessageId), zap.String("message type", msg.Type))
			_ = msg.Ack(false)
			return msg.Body, nil
		}
	}
	return nil, errors.New("err when getting notification status on channel")
}

func (c *rabbitmqClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *rabbitmqClient) configureQueue(queueName string) error {
	_, err := c.ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		utils.Log.Error("error on queue declare a amqp channel", zap.Error(err))
		return err
	}

	c.notificationStatus, err = c.ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		true,      // no-wait
		nil,       // args
	)
	return err
}
