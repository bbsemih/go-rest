package rabbitMQ

import (
	"context"
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueName = "transaction_queue"
)

type RabbitMQClient struct {
	conn              *amqp.Connection
	ch                *amqp.Channel
	connString        string
	transactionStatus <-chan amqp.Delivery
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NewRabbitMQClient(connString string) (*RabbitMQClient, error) {
	c := &RabbitMQClient{}
	var err error

	c.conn, err = amqp.Dial(connString)
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

func (c *RabbitMQClient) ConsumeByTransactionID(ctx context.Context, txID string) ([]byte, error) {
	for msg := range c.transactionStatus {
		if msg.MessageId == txID {
			_ = msg.Ack(false)
			return msg.Body, nil
		}
	}
	return nil, errors.New("failed to get transaction status on channel")
}

func (c *RabbitMQClient) Publish(ctx context.Context, txID int64) error {
	err := c.ch.Publish(
		"",
		QueueName,
		true,
		false,
		amqp.Publishing{
			MessageId: fmt.Sprintf("%d", txID),
			Body:      []byte(fmt.Sprintf("Transaction %d completed", txID)),
		})
	failOnError(err, "Failed to publish a message")

	return err
}

func (c *RabbitMQClient) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *RabbitMQClient) configureQueue() error {
	q, err := c.ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a rabbitMQ")

	err = c.ch.Qos(
		1, // the server will deliver that many messages to consumers before acks are received
		0,
		false,
	)
	failOnError(err, "Failed to set Qos")

	c.transactionStatus, err = c.ch.Consume(
		q.Name,
		"",
		false, // enabling this will cause the message to be deleted immediately after being received
		false,
		false,
		true,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	return err
}
