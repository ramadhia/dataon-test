package provider

import (
	"context"

	"github.com/ramadhia/bosnet/be/internal/config"
	"github.com/ramadhia/bosnet/be/internal/storage"
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
	db := storage.GetSqlServerDb()
	deferFn := func() {
		if db != nil {
			storage.CloseDB(db)
		}
	}

	provider.SetDB(db)

	return provider, deferFn, nil
}
