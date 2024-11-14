package model

import "time"

const (
	TOP_UP         = "TOP_UP"
	BANK_TRANSFER  = "BANK_TRANSFER"
	ONLINE_PAYMENT = "ONLINE_PAYMENT"

	TRX_TYPE_DEBIT  = "DEBIT"
	TRX_TYPE_CREDIT = "CREDIT"
)

type Transaction struct {
	ID                *string    `json:"id"`
	Status            *string    `json:"status"`
	UserID            *string    `json:"user_id"`
	TransactionMethod *string    `json:"transaction_method"`
	TransactionType   *string    `json:"transaction_type"`
	Amount            *float64   `json:"amount"`
	Remarks           *string    `json:"remarks"`
	BalanceBefore     *float64   `json:"balance_before"`
	BalanceAfter      *float64   `json:"balance_after"`
	CreatedDate       *time.Time `json:"created_date"`
	UpdatedDate       *time.Time `json:"updated_date"`
}

type Balance struct {
	UserID      *string    `json:"user_id"`
	Amount      *float64   `json:"amount"`
	LastUpdated *time.Time `json:"last_updated"`
}
