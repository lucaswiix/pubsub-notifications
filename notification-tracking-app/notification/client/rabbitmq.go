package client

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
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
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = c.configureQueue(queueName)

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
