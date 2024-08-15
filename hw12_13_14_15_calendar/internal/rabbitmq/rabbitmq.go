package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

// Client представляет интерфейс для работы с RabbitMQ.
type Client interface {
	Connect() error
	Close() error
	SendNotification(notification dto.NotificationData) error
	ReceiveNotifications(ctx context.Context) (<-chan dto.NotificationData, error)
}

// rabbitClient реализует интерфейс Client для работы с RabbitMQ.
type rabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	cfg     config.RabbitMQConfig
	logger  logger.Logger
}

// NewClient создает нового клиента для работы с RabbitMQ.
func NewClient(cfg config.RabbitMQConfig, log logger.Logger) (Client, error) {
	log.Infof("cfg= %v", cfg)

	return &rabbitClient{
		cfg:    cfg,
		logger: log,
	}, nil
}

func (c *rabbitClient) Connect() error {
	var err error

	c.conn, err = amqp.Dial(c.cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	c.logger.Info("connected to RabbitMQ")

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	c.logger.Infof("channel opened. name = %s", c.cfg.QueueName)

	c.logger.Infof("starting QueueDeclare. cfg.QueueName = %s", c.cfg.QueueName)

	c.queue, err = c.channel.QueueDeclare(
		c.cfg.QueueName, // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	c.logger.Infof("queue declared: name = %s", c.queue.Name)

	return nil
}

func (c *rabbitClient) Close() error {
	if err := c.channel.Close(); err != nil {
		return fmt.Errorf("on close channel: %w", err)
	}
	c.logger.Info("Channel closed")

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("on close connection: %w", err)
	}
	c.logger.Info("Connection closed")

	return nil
}

func (c *rabbitClient) SendNotification(notification dto.NotificationData) error {
	// Сериализация уведомления в JSON
	messageBody, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	c.logger.Infof("SendNotification c.queue.Name = %s", c.queue.Name)

	err = c.channel.Publish(
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		})
	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}
	c.logger.Info("Notification published: %s", notification.ID.String())
	return nil
}

func (c *rabbitClient) ReceiveNotifications(ctx context.Context) (<-chan dto.NotificationData, error) {
	c.logger.Infof("ReceiveNotifications c.queue.Name = %s", c.queue.Name)
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	notificationChannel := make(chan dto.NotificationData)
	go func() {
		defer close(notificationChannel) // Ensure the channel is closed when done

		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return // Channel closed, exit
				}
				var notification dto.NotificationData
				if err := json.Unmarshal(msg.Body, &notification); err != nil {
					c.logger.Error("Failed to unmarshal notification: %v", err)
					continue
				}
				notificationChannel <- notification
			case <-ctx.Done():
				// Context canceled, exit
				c.logger.Info("Context canceled, stopping notification receiver")
				return
			}
		}
	}()

	return notificationChannel, nil
}
