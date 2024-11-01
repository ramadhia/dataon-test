package repository

import (
	"context"
	"time"

	"github.com/ramadhia/bosnet/be/internal/model"
)

type HistoryRepository interface {
	FetchHistory(ctx context.Context, filter FetchHistoryFilter) ([]*model.History, error)
	AddHistory(ctx context.Context, filter AddHistoryFilter) (bool, error)
}

type FetchHistoryFilter struct {
	SzAccountId         *string
	StartDtmTransaction *string
	EndDtmTransaction   *string
	Offset              *int
	Limit               *int
}

type AddHistoryFilter struct {
	ToSzAccountId   string
	FromSzAccountId string
	SzCurrencyId    string
	DecAmount       float32
	SzNote          string
	DtmTransaction  time.Time
}
