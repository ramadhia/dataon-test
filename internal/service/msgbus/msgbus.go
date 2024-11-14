package msgbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ramadhia/mnc-test/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

type consumerConfig struct {
	topicName    string
	consumerName string
	handlerFunc  service.MessageHandler
	concurrency  int
}

// RabbitMQService adalah implementasi MessageBus menggunakan go-rabbitmq
type RabbitMQService struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
	consumers map[string]*rabbitmq.Consumer
	handlers  []consumerConfig
}

// NewRabbitMQService menginisialisasi RabbitMQService dengan publisher
func NewRabbitMQService(dsn string) *RabbitMQService {
	conn, err := rabbitmq.NewConn(dsn,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	publisher, err := rabbitmq.NewPublisher(conn, rabbitmq.WithPublisherOptionsLogging)
	if err != nil {
		return nil
	}

	return &RabbitMQService{
		conn:      conn,
		publisher: publisher,
		consumers: make(map[string]*rabbitmq.Consumer),
		handlers:  []consumerConfig{},
	}
}

// Publish mengirimkan pesan ke topik RabbitMQ yang ditentukan
func (r *RabbitMQService) Publish(ctx context.Context, topicName string, message *service.Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	headers := make(map[string]interface{})
	for key, value := range message.Headers {
		headers[key] = value
	}

	err = r.publisher.Publish(
		body,
		[]string{topicName},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsHeaders(headers),
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// On mendaftarkan handler untuk menerima pesan pada topik tertentu tanpa langsung menjalankan consumer
func (r *RabbitMQService) On(topicName string, consumerName string, handlerFunc service.MessageHandler, concurrency int) error {
	r.handlers = append(r.handlers, consumerConfig{
		topicName:    topicName,
		consumerName: consumerName,
		handlerFunc:  handlerFunc,
		concurrency:  concurrency,
	})
	return nil
}

func (r *RabbitMQService) StartConsuming() error {
	for _, handler := range r.handlers {
		//binding := fmt.Sprintf("%s:%s", handler.topicName, handler.consumerName)
		consumer, err := rabbitmq.NewConsumer(
			r.conn,
			handler.topicName,
			rabbitmq.WithConsumerOptionsConcurrency(handler.concurrency),
			rabbitmq.WithConsumerOptionsConsumerName(handler.consumerName),
			rabbitmq.WithConsumerOptionsQueueDurable,
			rabbitmq.WithConsumerOptionsLogging,
			rabbitmq.WithConsumerOptionsExchangeName(handler.topicName),
		)
		if err != nil {
			return fmt.Errorf("failed to create consumer: %w", err)
		}

		r.consumers[handler.topicName] = consumer

		go func(h consumerConfig) {
			err := consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
				confirm := &RabbitMQMessageConfirm{delivery: d}
				_ = h.handlerFunc(confirm, d.Body)
				return rabbitmq.Manual
			})
			if err != nil {
				fmt.Printf("error while consuming messages on topic %s: %v\n", h.topicName, err)
			}
		}(handler)
	}

	// Channel untuk menangani sinyal
	waitChan := make(chan os.Signal, 1)
	signal.Notify(waitChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Started consuming messages on all topics. Waiting for interrupt signal to shut down gracefully.")
	<-waitChan // Menunggu sinyal

	log.Println("Interrupt signal received. Closing RabbitMQ Service...")

	return r.Close()
}

func (r *RabbitMQService) Close() error {
	r.publisher.Close()

	for _, consumer := range r.consumers {
		consumer.Close()
	}
	return nil
}

// RabbitMQMessageConfirm implementasi MessageConfirm untuk go-rabbitmq
type RabbitMQMessageConfirm struct {
	delivery rabbitmq.Delivery
}

func (c *RabbitMQMessageConfirm) Ack() error {
	return c.delivery.Ack(false)
}

func (c *RabbitMQMessageConfirm) Nack() error {
	return c.delivery.Nack(false, true)
}

func (c *RabbitMQMessageConfirm) Retry(delay int64, maxRetry int) error {
	retryCount, _ := c.delivery.Headers["x-retry-count"].(int)
	if retryCount >= maxRetry {
		return fmt.Errorf("max retry reached")
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)
	return c.delivery.Nack(false, true)
}
