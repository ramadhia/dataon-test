package repository

import (
	"context"
	"github.com/ramadhia/bosnet/be/internal/model"
)

type BalanceRepository interface {
	GetBalance(ctx context.Context, szAccountId string, szCurrencyId *string) ([]*model.Balance, error)
}
