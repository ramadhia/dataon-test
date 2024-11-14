package transcation

import (
	"context"
	"errors"
	"github.com/ramadhia/mnc-test/internal/service"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/ramadhia/mnc-test/internal/config"
	"github.com/ramadhia/mnc-test/internal/model"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/repository"
	"github.com/ramadhia/mnc-test/internal/usecase"

	"github.com/samber/lo"
)

type TransactionImplement struct {
	config          config.Config
	transactionRepo repository.TransactionRepository
	balanceRepo     repository.BalanceRepository
	msgBus          service.MessageBus
}

func NewTransaction(p *provider.Provider) *TransactionImplement {
	return &TransactionImplement{
		config:          p.Config(),
		transactionRepo: p.TransactionRepo(),
		balanceRepo:     p.BalanceRepo(),
		msgBus:          p.MessageBus(),
	}
}

func (t *TransactionImplement) FetchTransaction(ctx context.Context, req usecase.FetchTransactionRequest) ([]*model.Transaction, error) {
	logger := logrus.WithField("method", "usecase.TransactionImplement.FetchTransaction")

	logger.Debugf("about to fetch the transactions")
	data, err := t.transactionRepo.FetchTransaction(ctx, repository.FetchTransactionFilter{
		UserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []*model.Transaction{}, nil
	}

	return data, nil
}

func (t *TransactionImplement) AddTransaction(ctx context.Context, claim model.Claim, args usecase.AddTransactionRequest) (*model.Transaction, error) {
	logger := logrus.WithField("method", "usecase.TransactionImplement.FetchTransaction")

	// Check balance / saldo before transaction
	userId := claim.ID

	var balanceCurrent, balanceAfter float64
	logger.Debugf("about to get the balance")
	balance, err := t.balanceRepo.GetBalance(ctx, userId)
	if err != nil {
		return nil, err
	}

	// checking balance first if trx type is DEBIT
	if args.GetTrxType() == model.TRX_TYPE_DEBIT {
		if balance == nil {
			return nil, errors.New("balance is not enough")
		}

		if *balance.Amount-*args.Amount < 1 {
			return nil, errors.New("balance is not enough")
		}
	}

	// get balance before
	balanceCurrent = t.getBalanceAmount(balance)

	if *args.Method == model.TOP_UP {
		balanceAfter = balanceCurrent + *args.Amount
	} else {
		balanceAfter = balanceCurrent - *args.Amount
	}

	logger.Debugf("about to add a transaction")
	// initialize transaction
	trx, err := t.transactionRepo.AddTransaction(ctx, model.Transaction{
		UserID:            &userId,
		Amount:            args.Amount,
		TransactionMethod: args.Method,
		TransactionType:   lo.ToPtr(args.GetTrxType()),
		BalanceBefore:     &balanceCurrent,
		BalanceAfter:      &balanceAfter,
		Remarks:           args.Remarks,
		Status:            lo.ToPtr("PENDING"),
		CreatedDate:       lo.ToPtr(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	// deduct balance when transaction type is debit
	if trx != nil && args.GetTrxType() == model.TRX_TYPE_DEBIT {

		_, err := t.balanceRepo.UpsertBalance(ctx, model.Balance{
			UserID:      trx.UserID,
			Amount:      lo.ToPtr(balanceCurrent - *args.Amount),
			LastUpdated: lo.ToPtr(time.Now()),
		})
		if err != nil {
			logger.Warning("error when update the balance after add transaction")
			return nil, err
		}
	}

	if trx != nil {
		message := service.NewMessage(service.Headers{
			"user_id": userId,
		}, usecase.WorkerTransactionRequest{
			RefersTrxID:       trx.ID,
			RefersUserID:      trx.UserID,
			TargetUserID:      args.TargetUser,
			TransactionMethod: trx.TransactionMethod,
			TransactionType:   trx.TransactionType,
			Amount:            args.Amount,
			Remarks:           args.Remarks,
		})

		logger.Debugf("about to publish the transfer data")
		if err := t.msgBus.Publish(ctx, "worker-transfer.consume", message); err != nil {
			logger.Warning("error when publishing the transfer data")
			return nil, err
		}
	}

	return trx, nil
}

func (t *TransactionImplement) WorkerTransaction(ctx context.Context, req usecase.WorkerTransactionRequest) error {
	logger := logrus.WithField("method", "usecase.TransactionImplement.WorkerTransaction")

	var balanceCurrent float64

	logger.Debugf("about to get the balance")
	balance, err := t.balanceRepo.GetBalance(ctx, *req.RefersUserID)
	if err != nil {
		logger.WithError(err).Warning("error when get the current balance")
		return err
	}

	// get balance before
	balanceCurrent = t.getBalanceAmount(balance)

	// update transaction when trx method is TOP-UP or Payment
	if *req.TransactionMethod == model.TOP_UP || *req.TransactionMethod == model.ONLINE_PAYMENT {

		// top up balance
		if *req.TransactionMethod == model.TOP_UP {

			_, err = t.balanceRepo.UpsertBalance(ctx, model.Balance{
				UserID:      balance.UserID,
				Amount:      lo.ToPtr(balanceCurrent + *req.Amount),
				LastUpdated: lo.ToPtr(time.Now()),
			})
			if err != nil {
				return err
			}
		}

		if err := t.transactionRepo.UpdateTransaction(ctx, model.Transaction{
			ID:          req.RefersTrxID,
			Status:      lo.ToPtr("SUCCESS"),
			UpdatedDate: lo.ToPtr(time.Now()),
		}); err != nil {
			return err
		}

		return nil
	}

	// checking balance first if trx method is TRANSFER
	if *req.TransactionMethod == model.BANK_TRANSFER {
		if balance == nil {
			logger.Warning("balance is not found")
			return nil
		}

		if balanceCurrent-*req.Amount < 1 {
			logger.Warning("balance is not enough")

			// update transaction status if balance is not enough
			if err := t.transactionRepo.UpdateTransaction(ctx, model.Transaction{
				ID:          req.RefersTrxID,
				Status:      lo.ToPtr("FAILED"),
				Remarks:     lo.ToPtr("Balance is not enough"),
				UpdatedDate: lo.ToPtr(time.Now()),
			}); err != nil {
				logger.Error("error when update the transactions")
				return err
			}

			return nil
		}
	}

	logger.Debugf("about to get the target balance")
	targetBalance, err := t.balanceRepo.GetBalance(ctx, *req.TargetUserID)
	if err != nil {
		logger.WithError(err).Warning("error when get the target balance")
		return err
	}

	// initialize transaction when trx method is TRANSFER
	trx, err := t.transactionRepo.AddTransaction(ctx, model.Transaction{
		UserID:            req.TargetUserID,
		Amount:            req.Amount,
		TransactionMethod: req.TransactionMethod,
		TransactionType:   lo.ToPtr(model.TRX_TYPE_CREDIT),
		BalanceBefore:     targetBalance.Amount,
		BalanceAfter:      lo.ToPtr(*targetBalance.Amount + *req.Amount),
		Remarks:           req.Remarks,
		Status:            lo.ToPtr("SUCCESS"),
		CreatedDate:       lo.ToPtr(time.Now()),
	})
	if err != nil {
		logger.WithError(err).Warning("error when add a transaction for target user")
		return err
	}

	if trx != nil {
		if err := t.processAfterTransfer(ctx, req, targetBalance); err != nil {
			logger.WithError(err).Warning("error when process after payment transfer")
			return err
		}
	}

	logger.Info("All transactions already processed")
	return nil
}

func (t *TransactionImplement) getBalanceAmount(balance *model.Balance) (b float64) {
	if balance != nil && balance.Amount != nil {
		b = *balance.Amount
	}
	return
}

func (t *TransactionImplement) processAfterTransfer(ctx context.Context, req usecase.WorkerTransactionRequest, targetBalance *model.Balance) (err error) {
	// update the refers transaction status and balance if all transactions already processed
	if err := t.transactionRepo.UpdateTransaction(ctx, model.Transaction{
		ID:          req.RefersTrxID,
		Status:      lo.ToPtr("SUCCESS"),
		UpdatedDate: lo.ToPtr(time.Now()),
	}); err != nil {
		return err
	}

	// get target balance and update amount balance
	_, err = t.balanceRepo.UpsertBalance(ctx, model.Balance{
		UserID:      targetBalance.UserID,
		Amount:      lo.ToPtr(*targetBalance.Amount + *req.Amount),
		LastUpdated: lo.ToPtr(time.Now()),
	})
	if err != nil {
		return err
	}

	return
}
