package messaging

import (
	"github.com/ramadhia/mnc-test/internal/handler/messaging/handler"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/service"
)

const (
	DefaultConsumerName = "worker.app"
)

func RegisterSubscribers(p *provider.Provider, topics []string) error {
	// set messaging controller

	transfer := handler.NewTransfer(p)

	topicHandlers := map[string][]struct {
		consumerName string
		handler      service.MessageHandler
	}{
		// subscribe topic
		"worker-transfer.consume": {{DefaultConsumerName, transfer.Transfer}},
	}

	for topic, handlers := range topicHandlers {
		for _, h := range handlers {
			err := p.MessageBus().On(topic, h.consumerName, h.handler, 1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
