package gorm

import (
	"context"
	"errors"

	"github.com/ramadhia/dataon-test/internal/entity"
	"github.com/ramadhia/dataon-test/internal/repository"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationGorm struct {
	db *gorm.DB
}

func NewOrganizationRepo(db *gorm.DB) *OrganizationGorm {
	return &OrganizationGorm{
		db: db,
	}
}

func (t *OrganizationGorm) FetchOrganization(ctx context.Context, filter repository.FetchTransactionFilter) (ret []*entity.Organization, err error) {
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
	var items []Organization
	if err = q.Find(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var transaction Organization
	return transaction.ToModels(items), nil
}

func (t *OrganizationGorm) FetchComplete(ctx context.Context, filter repository.FetchTransactionFilter) ([]*entity.Organization, error) {
	logger := logrus.WithField("method", "repo.TransactionGorm.FetchOrganization")

	logger.Debugf("about to fetch the complete transactions")

	fetch, err := t.organizationRecursive(ctx, nil, 1, 5)
	if err != nil {
		return nil, err
	}

	org := Organization{}
	return org.ToModels(fetch), nil
}

func (t *OrganizationGorm) AddOrganization(ctx context.Context, data entity.Organization) (ret *entity.Organization, err error) {
	logger := logrus.WithField("method", "repo.TransactionGorm.AddTransaction")
	q := t.db.WithContext(ctx)

	logger.Debugf("about to add a transactions")
	item := Organization{}
	trx := item.FromModel(data)
	if err := q.Create(&trx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return trx.ToModel(), nil
}

func (t *OrganizationGorm) UpdateOrganization(ctx context.Context, data entity.Organization) (err error) {
	logger := logrus.WithField("method", "repo.TransactionGorm.UpdateTransaction")

	tx := t.db.WithContext(ctx)

	var item Organization
	trx := item.FromModel(data)

	logger.Debugf("about to update transaction data")
	err = tx.Model(Organization{
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

func (t *OrganizationGorm) organizationRecursive(ctx context.Context, parentID *string, level int, maxLevel int) ([]Organization, error) {
	tx := t.db.WithContext(ctx)
	if level > maxLevel {
		return nil, nil
	}

	var organizations []Organization
	q := tx.Preload("Group").
		Preload("Group.Users").
		Preload("Children.Group").
		Preload("Children.Group.Users")

	if parentID != nil {
		q = q.Where("parent_id = ?", parentID)
	} else {
		q = q.Where("parent_id IS NULL")
	}

	err := q.Find(&organizations).Error
	if err != nil {
		return nil, err
	}

	for i := range organizations {
		children, err := t.organizationRecursive(ctx, organizations[i].ID, level+1, maxLevel)
		if err != nil {
			return nil, err
		}
		organizations[i].Children = children
	}

	return organizations, nil
}

type Organization struct {
	ID       *string        `gorm:"primaryKey;type:varchar(50)" json:"id"`
	ParentID *string        `gorm:"type:varchar(20);" json:"parent_id"`
	GroupID  *string        `gorm:"type:varchar(20);not null" json:"group_id"`
	Group    *Group         `gorm:"foreignKey:GroupID;references:ID"` // Relasi ke Group
	Children []Organization `gorm:"foreignKey:ParentID"`
}

func (t *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == nil {
		t.ID = lo.ToPtr(uuid.New().String())

	}
	return
}

func (t *Organization) FromModel(e entity.Organization) *Organization {
	return &Organization{
		ID:       e.ID,
		ParentID: e.ParentID,
		GroupID:  e.GroupID,
	}
}

func (t *Organization) ToModel() *entity.Organization {
	var group *entity.Group
	if t.Group != nil {
		group = t.Group.ToModel()
	}

	var children []entity.Organization
	if len(t.Children) > 0 {
		for _, v := range t.Children {
			children = append(children, *v.ToModel())
		}
	}

	m := &entity.Organization{
		ID:       t.ID,
		ParentID: t.ParentID,
		GroupID:  t.GroupID,
		Group:    group,
		Children: children,
	}

	return m
}

func (t *Organization) ToModels(d []Organization) []*entity.Organization {
	var models []*entity.Organization
	for _, v := range d {
		models = append(models, v.ToModel())
	}
	return models
}
