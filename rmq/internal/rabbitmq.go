package internal

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// rule of thumb is to use a single connection per app and spawn channels for every task
	conn *amqp.Connection // a tcp connection used by the client
	ch   *amqp.Channel    // a multiplexed connection over the tcp connection i.e, conn
}

func ConnectRabbitMQ(username string, password string,
	host string, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}
	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) NewExchangeDeclare(exchangeName, kind string, durable, autodelete bool) error {
	err := rc.ch.ExchangeDeclare(exchangeName, kind, durable, autodelete, false, false, nil)
	return err
}

func (rc RabbitClient) NewQueueDeclare(queueName string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}

func (rc RabbitClient) CreateBinding(name string, binding string, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	return rc.ch.PublishWithContext(ctx, exchange, routingKey, true, false, options)
}

func (rc RabbitClient) Receive(ctx context.Context, queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.ConsumeWithContext(ctx, queue, consumer, autoAck, false, false, false, nil)
}
