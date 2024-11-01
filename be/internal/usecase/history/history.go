package history

import (
	"context"
	"errors"
	"time"

	"github.com/ramadhia/bosnet/be/internal/config"
	"github.com/ramadhia/bosnet/be/internal/model"
	"github.com/ramadhia/bosnet/be/internal/provider"
	"github.com/ramadhia/bosnet/be/internal/repository"
	"github.com/ramadhia/bosnet/be/internal/usecase"

	"github.com/shortlyst-ai/go-helper"
)

type HistoryImpl struct {
	config      config.Config
	historyRepo repository.HistoryRepository
	balanceRepo repository.BalanceRepository
}

func NewHistory(p *provider.Provider) *HistoryImpl {
	return &HistoryImpl{
		config:      p.Config(),
		historyRepo: p.HistoryRepo(),
		balanceRepo: p.BalanceRepo(),
	}
}

func (t HistoryImpl) FetchHistory(ctx context.Context, request usecase.HistoryRequest) (*usecase.FetchHistoryResponse, error) {
	data, err := t.historyRepo.FetchHistory(ctx, repository.FetchHistoryFilter{
		SzAccountId:         request.SzAccountId,
		StartDtmTransaction: request.StartDtmTransaction,
		EndDtmTransaction:   request.EndDtmTransaction,
		Offset:              helper.Pointer(1),
		Limit:               helper.Pointer(5),
	})
	if err != nil {
		return nil, err
	}

	return &usecase.FetchHistoryResponse{
		Data: data,
	}, nil
}

func (t HistoryImpl) AddHistory(ctx context.Context, args usecase.AddHistoryRequest) (*usecase.FetchHistoryResponse, error) {
	var fromAccount, toAccount string
	if args.FromSzAccountId != nil {
		fromAccount = *args.FromSzAccountId
	}

	if args.ToSzAccountId != nil {
		toAccount = *args.ToSzAccountId
	}

	var data interface{}

	// Check balance jika bukan SETOR
	if args.SzNote != nil && *args.SzNote != "SETOR" {
		res, err := t.balanceRepo.GetBalance(ctx, *args.FromSzAccountId, args.SzCurrencyId)
		if err != nil {
			return nil, err
		}

		// return nil if account not found
		if res == nil {
			return &usecase.FetchHistoryResponse{
				Data: data,
			}, nil
		}
		var balance model.Balance
		balance = *res[0]

		if *balance.DecAmount-*args.DecAmount < 1 {
			return nil, errors.New("not enough balance")
		}

		data = res
	}

	if args.SzNote != nil && *args.SzNote == "TARIK" {
		args.DecAmount = helper.Pointer(0 - *args.DecAmount)
		toAccount = fromAccount
	}

	data, err := t.historyRepo.AddHistory(ctx, repository.AddHistoryFilter{
		ToSzAccountId:   toAccount,
		FromSzAccountId: fromAccount,
		SzCurrencyId:    *args.SzCurrencyId,
		DecAmount:       *args.DecAmount,
		SzNote:          *args.SzNote,
		DtmTransaction:  time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &usecase.FetchHistoryResponse{
		Data: data,
	}, nil
}
