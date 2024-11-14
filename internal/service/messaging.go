package service

import "context"

type MessageBus interface {
	Publish(ctx context.Context, topicName string, message *Message) error
	On(topicName string, consumerName string, handlerFunc MessageHandler, concurrency int) error
	Close() error
	StartConsuming() error
}

type MessageHandler func(MessageConfirm, []byte) error

type MessageConfirm interface {
	Ack() error
	Nack() error
	Retry(delay int64, maxRetry int) error
}

type Headers map[string]string

type Message struct {
	Headers Headers     `json:"headers"`
	Body    interface{} `json:"body"`
}

func NewMessage(headers Headers, body interface{}) *Message {
	return &Message{Headers: headers, Body: body}
}
