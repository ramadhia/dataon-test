package usecase

import (
	"context"
)

type HistoryUsecase interface {
	FetchHistory(ctx context.Context, args HistoryRequest) (*FetchHistoryResponse, error)
	AddHistory(ctx context.Context, args AddHistoryRequest) (*FetchHistoryResponse, error)
}

type HistoryRequest struct {
	SzAccountId         *string `form:"szAccountId"`
	StartDtmTransaction *string `form:"startDtmTransaction"`
	EndDtmTransaction   *string `form:"endDtmTransaction"`
}

type FetchHistoryResponse struct {
	Data interface{} `json:"data"`
}

type AddHistoryRequest struct {
	ToSzAccountId   *string  `json:"toSzAccountId"`
	FromSzAccountId *string  `json:"fromSzAccountId"`
	SzCurrencyId    *string  `json:"szCurrencyId"`
	DecAmount       *float32 `json:"decAmount"`
	SzNote          *string  `json:"szNote"`
}
