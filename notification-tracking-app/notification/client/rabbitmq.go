package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucaswiix/notifications-tracking-app/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type rabbitmqClient struct {
	conn               *amqp.Connection
	ch                 *amqp.Channel
	connString         string
	notificationStatus <-chan amqp.Delivery
	done               chan error
	lastRecoverTime    int64
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
			utils.Log.Debug(fmt.Sprintf("receive notification to user %s, type %s", msg.MessageId, msg.Type))
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

	go func() {
		utils.Log.Info("closing: ", zap.Error(<-c.conn.NotifyClose(make(chan *amqp.Error))))
		c.done <- errors.New("channel closed")
	}()

	c.notificationStatus, err = c.ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		true,      // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	return err
}
