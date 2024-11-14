package repository

import (
	"context"
	"github.com/ramadhia/mnc-test/internal/model"
)

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID string) (*model.Balance, error)
	UpsertBalance(ctx context.Context, balance model.Balance) (*model.Balance, error)
}
