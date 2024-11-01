package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/ramadhia/bosnet/be/internal/model"
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

func (p *BalanceGorm) GetBalance(ctx context.Context, szAccountId string, szCurrencyId *string) (ret []*model.Balance, err error) {
	q := p.db.WithContext(ctx)

	var items []BosBalance
	query := fmt.Sprintf("SELECT szAccountId, szCurrencyId, decAmount FROM BOS_Balance WHERE szAccountId = '%s'", szAccountId)
	if szCurrencyId != nil {
		query = fmt.Sprintf("SELECT szAccountId, szCurrencyId, decAmount FROM BOS_Balance WHERE szAccountId = '%s' AND szCurrencyId = '%s'", szAccountId, *szCurrencyId)
	}

	if err = q.Raw(query).Scan(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return BosBalance{}.ToModels(items), nil
}

type BosBalance struct {
	SzAccountId  *string  `gorm:"column:szAccountId"`
	SzCurrencyId *string  `gorm:"column:szCurrencyId"`
	DecAmount    *float32 `gorm:"column:decAmount"`
}

func (BosBalance) FromModel(model model.Balance) *BosBalance {
	return &BosBalance{
		SzAccountId:  model.SzAccountId,
		SzCurrencyId: model.SzCurrencyId,
		DecAmount:    model.DecAmount,
	}
}

func (BosBalance) TableName() string {
	return "BOS_Balance"
}

func (b BosBalance) ToModel() *model.Balance {
	m := &model.Balance{
		SzAccountId:  b.SzAccountId,
		SzCurrencyId: b.SzCurrencyId,
		DecAmount:    b.DecAmount,
	}

	return m
}

func (b BosBalance) ToModels(d []BosBalance) []*model.Balance {
	var models []*model.Balance
	for _, v := range d {
		models = append(models, v.ToModel())
	}
	return models
}
