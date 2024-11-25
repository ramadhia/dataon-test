package repository

import (
	"context"
	"time"

	"github.com/ramadhia/dataon-test/internal/entity"
)

type OrganizationRepository interface {
	FetchOrganization(ctx context.Context, filter FetchTransactionFilter) ([]*entity.Organization, error)
	FetchComplete(ctx context.Context, filter FetchTransactionFilter) ([]*entity.Organization, error)
	AddOrganization(ctx context.Context, transaction entity.Organization) (*entity.Organization, error)
	UpdateOrganization(ctx context.Context, transaction entity.Organization) error
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
