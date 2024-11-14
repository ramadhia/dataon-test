package repository

import (
	"context"
	"time"

	"github.com/ramadhia/mnc-test/internal/model"
)

type TransactionRepository interface {
	FetchTransaction(ctx context.Context, filter FetchTransactionFilter) ([]*model.Transaction, error)
	AddTransaction(ctx context.Context, transaction model.Transaction) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction model.Transaction) error
}

type FetchTransactionFilter struct {
	UserID      *string
	Offset      *int
	Limit       *int
	CreatedDate *time.Time
}

type AddTransactionFilter struct {
	FromUserID    *string    `json:"from_user_id"`
	ToUserID      *string    `json:"to_user_id"`
	Method        *string    `json:"method"`
	Amount        *float32   `json:"amount"`
	Type          *string    `json:"type"`
	BalanceBefore *float32   `json:"balance_before"`
	BalanceAfter  *float32   `json:"balance_after"`
	Remarks       *string    `json:"remarks"`
	Status        *string    `json:"status"`
	CreatedDate   *time.Time `json:"created_date"`
}
