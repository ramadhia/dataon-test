package response

import "time"

type TopupResponse struct {
	TopUpID       *string    `json:"top_up_id"`
	AmountTopUp   *float64   `json:"amount_top_up"`
	BalanceBefore *float64   `json:"balance_before"`
	BalanceAfter  *float64   `json:"balance_after"`
	CreatedDate   *time.Time `json:"created_date"`
}

type PaymentResponse struct {
	PaymentID     *string    `json:"payment_id"`
	Amount        *float64   `json:"amount"`
	Remarks       *string    `json:"remarks"`
	BalanceBefore *float64   `json:"balance_before"`
	BalanceAfter  *float64   `json:"balance_after"`
	CreatedDate   *time.Time `json:"created_date"`
}

type TransferResponse struct {
	TransferID    *string    `json:"transfer_id"`
	TargetUser    *string    `json:"target_user"`
	Amount        *float64   `json:"amount"`
	Remarks       *string    `json:"remarks"`
	BalanceBefore *float64   `json:"balance_before"`
	BalanceAfter  *float64   `json:"balance_after"`
	CreatedDate   *time.Time `json:"created_date"`
}
