package handler

import (
	"context"
	"encoding/json"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/service"
	"github.com/ramadhia/mnc-test/internal/usecase"
	"github.com/shortlyst-ai/go-helper"
	"github.com/sirupsen/logrus"
)

type Transfer struct {
	provider *provider.Provider
}

func NewTransfer(appContainer *provider.Provider) *Transfer {
	if appContainer == nil {
		panic("nil container")
	}
	return &Transfer{provider: appContainer}
}

func (t *Transfer) Transfer(confirm service.MessageConfirm, message []byte) error {
	logger := logrus.WithField("method", "message.Transfer.Transfer")

	logger.Debugf("sudah masuk sini")

	headers := make(map[string]string)
	req := usecase.WorkerTransactionRequest{}
	msgObj := service.Message{Headers: headers, Body: &req}
	if err := json.Unmarshal(message, &msgObj); err != nil {
		logger.WithError(err).Warning("Failed unmarshalling message")
		return confirm.Ack()
	}

	logger.Debugf(helper.MustJsonString(req))

	err := t.provider.TransactionUseCase().WorkerTransaction(context.Background(), req)
	if err != nil {
		logger.WithError(err).Error("Failed in transaction")
		return confirm.Nack()
	}

	return confirm.Ack()
}
