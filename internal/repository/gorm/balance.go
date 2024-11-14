package gorm

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/ramadhia/mnc-test/internal/model"

	"gorm.io/gorm"
)

type BalanceGorm struct {
	db *gorm.DB
}

func NewBalanceRepo(db *gorm.DB) *BalanceGorm {
	return &BalanceGorm{
		db: db,
	}
}

func (b *BalanceGorm) GetBalance(ctx context.Context, userID string) (ret *model.Balance, err error) {
	logger := logrus.WithField("method", "repo.BalanceGorm.GetBalance")
	q := b.db.WithContext(ctx)

	item := Balance{}
	q = q.Where("user_id = ?", userID)

	logger.Debugf("about to get data balance")
	if err = q.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, err
	}

	return item.ToModel(), nil
}

func (b *BalanceGorm) UpsertBalance(ctx context.Context, data model.Balance) (res *model.Balance, err error) {
	logger := logrus.WithField("method", "repo.BalanceGorm.UpsertBalance")

	tx := b.db.WithContext(ctx)

	logger.Debugf("about to get current balance")
	res, err = b.GetBalance(ctx, *data.UserID)
	if err != nil {
		return nil, err
	}

	item := Balance{}
	balance := item.FromModel(data)

	// insert new data if balance not found
	if res == nil {
		logger.Debugf("about to add a balance data")
		if err := tx.Create(&balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}
			return nil, err
		}

		return item.ToModel(), nil
	}

	// update if user_id is found
	logger.Debugf("about to update balance data")
	err = tx.Model(Balance{}).
		Where("user_id = ?", data.UserID).
		Updates(balance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		return nil, err
	}

	return
}

type Balance struct {
	UserID      *string
	Amount      *float64
	LastUpdated *time.Time
}

func (Balance) FromModel(model model.Balance) *Balance {
	return &Balance{
		UserID:      model.UserID,
		Amount:      model.Amount,
		LastUpdated: model.LastUpdated,
	}
}

func (Balance) TableName() string {
	return "balances"
}

func (b Balance) ToModel() *model.Balance {
	m := &model.Balance{
		UserID:      b.UserID,
		Amount:      b.Amount,
		LastUpdated: b.LastUpdated,
	}

	return m
}

func (b Balance) ToModels(d []Balance) []*model.Balance {
	var models []*model.Balance
	for _, v := range d {
		models = append(models, v.ToModel())
	}
	return models
}
