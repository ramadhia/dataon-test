package gorm

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/ramadhia/mnc-test/internal/model"
	"github.com/ramadhia/mnc-test/internal/repository"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type TransactionGorm struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionGorm {
	return &TransactionGorm{
		db: db,
	}
}

func (t *TransactionGorm) FetchTransaction(ctx context.Context, filter repository.FetchTransactionFilter) (ret []*model.Transaction, err error) {
	logger := logrus.WithField("method", "repo.TransactionGorm.FetchTransaction")

	q := t.db.WithContext(ctx)

	if filter.Limit != nil && filter.Offset != nil {
		limit := 5
		if filter.Limit != nil {
			limit = *filter.Limit
		}

		offset := 0
		if filter.Offset != nil {
			offset = *filter.Offset
		}

		q = q.Limit(limit).Offset(offset)
	}

	if filter.UserID != nil {
		q = q.Where("user_id = ?", filter.UserID)
	}

	if v := filter.CreatedDate; v != nil {
		q = q.Where("created_date >= ?", *v)
	}

	logger.Debugf("about to fetch the transactions")
	var items []Transaction
	if err = q.Find(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var transaction Transaction
	return transaction.ToModels(items), nil
}

func (t *TransactionGorm) AddTransaction(ctx context.Context, data model.Transaction) (ret *model.Transaction, err error) {
	logger := logrus.WithField("method", "repo.TransactionGorm.AddTransaction")
	q := t.db.WithContext(ctx)

	logger.Debugf("about to add a transactions")
	item := Transaction{}
	trx := item.FromModel(data)
	if err := q.Create(&trx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return trx.ToModel(), nil
}

func (t *TransactionGorm) UpdateTransaction(ctx context.Context, data model.Transaction) (err error) {
	logger := logrus.WithField("method", "repo.TransactionGorm.UpdateTransaction")

	tx := t.db.WithContext(ctx)

	var item Transaction
	trx := item.FromModel(data)

	logger.Debugf("about to update transaction data")
	err = tx.Model(Transaction{
		ID: data.ID,
	}).Updates(trx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		return err
	}

	return
}

type Transaction struct {
	ID                *string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Status            *string    `gorm:"type:varchar(20);not null" json:"status"`
	UserID            *string    `gorm:"type:varchar(50);not null" json:"user_id"`
	TransactionMethod *string    `gorm:"type:varchar(50);not null" json:"transaction_method"`
	TransactionType   *string    `gorm:"type:varchar(10);not null" json:"transaction_type"`
	Amount            *float64   `gorm:"type:numeric(15,2);not null" json:"amount"`
	Remarks           *string    `gorm:"type:text" json:"remarks"`
	BalanceBefore     *float64   `gorm:"type:numeric(15,2);not null" json:"balance_before"`
	BalanceAfter      *float64   `gorm:"type:numeric(15,2);not null" json:"balance_after"`
	CreatedDate       *time.Time `gorm:"type:timestamptz;not null" json:"created_date"`
	UpdatedDate       *time.Time `gorm:"type:timestamptz;not null" json:"updated_date"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == nil {
		t.ID = lo.ToPtr(uuid.New().String())

	}
	return
}

func (t *Transaction) FromModel(model model.Transaction) *Transaction {
	return &Transaction{
		ID:                model.ID,
		Status:            model.Status,
		UserID:            model.UserID,
		TransactionMethod: model.TransactionMethod,
		TransactionType:   model.TransactionType,
		Amount:            model.Amount,
		Remarks:           model.Remarks,
		BalanceBefore:     model.BalanceBefore,
		BalanceAfter:      model.BalanceAfter,
		CreatedDate:       model.CreatedDate,
		UpdatedDate:       model.UpdatedDate,
	}
}

func (t *Transaction) ToModel() *model.Transaction {
	m := &model.Transaction{
		ID:                t.ID,
		Status:            t.Status,
		UserID:            t.UserID,
		TransactionMethod: t.TransactionMethod,
		TransactionType:   t.TransactionType,
		Amount:            t.Amount,
		Remarks:           t.Remarks,
		BalanceBefore:     t.BalanceBefore,
		BalanceAfter:      t.BalanceAfter,
		CreatedDate:       t.CreatedDate,
	}

	return m
}

func (t *Transaction) ToModels(d []Transaction) []*model.Transaction {
	var models []*model.Transaction
	for _, v := range d {
		models = append(models, v.ToModel())
	}
	return models
}
