package model

import "time"

type History struct {
	SzTransactionId *string    `json:"szTransactionId"`
	SzAccountId     *string    `json:"szAccountId"`
	SzCurrencyId    *string    `json:"szCurrencyId"`
	DecAmount       *float32   `json:"decAmount"`
	SzNote          *string    `json:"szNote"`
	DtmTransaction  *time.Time `json:"dtmTransaction"`
}

type Balance struct {
	SzAccountId  *string  `json:"szAccountId"`
	SzCurrencyId *string  `json:"szCurrencyId"`
	DecAmount    *float32 `json:"decAmount"`
}

type Counter struct {
	SzCurrencyId string `json:"szCurrencyId"`
	ILastNumber  string `json:"iLastNumber"`
}
