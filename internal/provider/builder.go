package provider

import (
	"context"
	"fmt"
	"github.com/ramadhia/dataon-test/internal/config"
	"github.com/ramadhia/dataon-test/internal/service/msgbus"
	"github.com/ramadhia/dataon-test/internal/storage"
)

type ProviderBuilder interface {
	// Build will return container instance, which contains dependencies and resource clean up function
	// will return nil, nil, error when error happen
	Build(ctx context.Context) (*Provider, func(), error)
}

type DefaultProviderBuilder struct {
}

func (DefaultProviderBuilder) Build(ctx context.Context) (*Provider, func(), error) {
	cfg := config.Instance()

	// init provider
	provider := NewProvider()
	provider.SetConfig(cfg)

	// init db
	db := storage.GetPostgresDb()

	// init message bus
	msgBusSvc := initializeMessaging(1)
	deferFn := func() {
		if db != nil {
			storage.CloseDB(db)
		}
		if msgBusSvc != nil {
			msgBusSvc.Close()
		}
	}

	provider.SetDB(db)
	provider.SetMessageBus(msgBusSvc)

	return provider, deferFn, nil
}

func initializeMessaging(threadNum int) *msgbus.RabbitMQService {
	cfg := config.Instance() // check env
	if cfg.RabbitMq.Host == "" {
		panic("No 'amqp_server_url' set in configuration, cannot start")
	}
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s%s", cfg.RabbitMq.User, cfg.RabbitMq.Password, cfg.RabbitMq.Host, cfg.RabbitMq.Port, "/")

	mc := msgbus.NewRabbitMQService(dsn)
	return mc
}
