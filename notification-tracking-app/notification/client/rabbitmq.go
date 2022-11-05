package client

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	QueueName = DefaultRabbitConfig.QueueName
)

type rabbitmqClient struct {
	conn               *amqp.Connection
	ch                 *amqp.Channel
	connString         string
	notificationStatus <-chan amqp.Delivery
}

func NewRabbitMQClient(connectionString string) (*rabbitmqClient, error) {
	c := &rabbitmqClient{}
	var err error

	c.conn, err = amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.configureQueue()

	return c, err
}

func (c *rabbitmqClient) ConsumeByUserID(ctx context.Context, userID, nType string) ([]byte, error) {
	for msg := range c.notificationStatus {
		if msg.MessageId == userID && msg.Type == nType {
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

func (c *rabbitmqClient) configureQueue() error {
	_, err := c.ch.QueueDeclare(
		QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	c.notificationStatus, err = c.ch.Consume(
		QueueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		true,      // no-wait
		nil,       // args
	)
	return err
}
