package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/eventscompass/service-framework/service"
)

// AMQPBus is a message bus backed by RabbitMQ message broker.
type AMQPBus struct {
	// conn is the connection to the RabbitMQ message broker.
	conn *amqp.Connection

	// exchange is the exchange associated with this Bus.
	exchange string
}

var (
	_ service.MessageBus = (*AMQPBus)(nil)
	_ io.Closer          = (*AMQPBus)(nil)
)

// NewAMQPBus creates a new [AMQPBus] instance.
func NewAMQPBus(cfg *service.BusConfig, exchange string) (*AMQPBus, error) {
	connInfo := fmt.Sprintf(
		"amqp://%s:%s@%s:%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(connInfo)
	if err != nil { // maybe use exponential backoff for connecting ?
		return nil, fmt.Errorf("%w: amqp dial: %v", service.ErrUnexpected, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%w: pubsub channel: %v", service.ErrUnexpected, err)
	}
	defer ch.Close()

	return &AMQPBus{
		conn:     conn,
		exchange: exchange,
	}, nil
}

// Publish implements the [service.MessageBus] interface.
func (b *AMQPBus) Publish(ctx context.Context, p service.Payload) error {
	if b.conn.IsClosed() {
		return service.ErrConnectionClosed
	}

	topic := p.Topic()
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("%w: marshal payload: %v", service.ErrUnexpected, err)
	}

	// Note that AMQP channels are not thread-safe. Thus, we will be creating a
	// new channel for every published message. By using separate AMQP channels
	// we can reuse the same AMQP connection.
	ch, err := b.conn.Channel()
	if err != nil {
		return fmt.Errorf("%w: pubsub channel: %v", service.ErrUnexpected, err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(b.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("%w: exchange declare: %v", service.ErrUnexpected, err)
	}

	return ch.PublishWithContext(
		ctx,
		b.exchange, // exchange
		topic,      // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Subscribe implements the [service.MessageBus] interface.
func (b *AMQPBus) Subscribe(ctx context.Context, topic string, h service.EventHandler) error {
	if b.conn.IsClosed() {
		return service.ErrConnectionClosed
	}

	ch, err := b.conn.Channel()
	if err != nil {
		return fmt.Errorf("%w: pubsub channel: %v", service.ErrUnexpected, err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("%w: queue declare: %v", service.ErrUnexpected, err)
	}
	err = ch.QueueBind(q.Name, topic, b.exchange, false, nil)
	if err != nil {
		return fmt.Errorf("%w: queue bind: %v", service.ErrUnexpected, err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // non-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("%w: queue consume: %v", service.ErrUnexpected, err)
	}

	for msg := range msgs {
		var payload any
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			return fmt.Errorf("%w: unmarshal payload: %v", service.ErrUnexpected, err)
		}
		if err := h(ctx, payload); err != nil {
			return fmt.Errorf("%w: event handler: %v", service.ErrUnexpected, err)
		}
		_ = msg.Ack(false)
	}

	return nil
}

// Close implements the [io.Closer] interface.
func (b *AMQPBus) Close() error {
	return b.conn.Close()
}