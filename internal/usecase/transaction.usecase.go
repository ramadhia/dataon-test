package usecase

import (
	"context"
	"time"

	"github.com/ramadhia/mnc-test/internal/model"

	"github.com/go-playground/validator"
	"github.com/shortlyst-ai/go-helper"
)

type TransactionUsecase interface {
	FetchTransaction(ctx context.Context, args FetchTransactionRequest) ([]*model.Transaction, error)
	AddTransaction(ctx context.Context, claim model.Claim, args AddTransactionRequest) (*model.Transaction, error)
	WorkerTransaction(ctx context.Context, req WorkerTransactionRequest) error
}

type FetchTransactionRequest struct {
	UserID *string
}

type FetchTransacitonResponse struct {
	Result interface{} `json:"result"`
}

type AddTransactionRequest struct {
	TargetUser *string  `json:"target_user,omitempty"`
	Method     *string  `json:"method,omitempty"`
	Amount     *float64 `json:"amount,omitempty" validate:"required,numeric"`
	Remarks    *string  `json:"remarks,omitempty"`
}

func (a *AddTransactionRequest) Validate() (err error) {
	validate = validator.New()
	if err := validate.Struct(a); err != nil {
		errMsg := err.Error()
		return helper.NewParameterError(&errMsg)
	}
	return
}

func (a *AddTransactionRequest) GetTrxType() (trxType string) {
	if *a.Method == model.TOP_UP {
		trxType = model.TRX_TYPE_CREDIT
	} else {
		trxType = model.TRX_TYPE_DEBIT
	}
	return
}

type WorkerTransactionRequest struct {
	RefersTrxID       *string    `json:"refers_trx_id,omitempty"`
	RefersUserID      *string    `json:"refers_user_id,omitempty"`
	TargetUserID      *string    `json:"target_user_id,omitempty"`
	Status            *string    `json:"status"`
	TransactionMethod *string    `json:"transaction_method"`
	TransactionType   *string    `json:"transaction_type"`
	Amount            *float64   `json:"amount"`
	Remarks           *string    `json:"remarks"`
	CreatedDate       *time.Time `json:"created_date"`
}

func (a *WorkerTransactionRequest) GetTrxType() (trxType string) {
	if *a.TransactionMethod == model.TOP_UP {
		trxType = model.TRX_TYPE_CREDIT
	} else {
		trxType = model.TRX_TYPE_DEBIT
	}
	return
}
