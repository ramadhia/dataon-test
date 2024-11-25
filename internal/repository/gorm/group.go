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

type GroupGorm struct {
	db *gorm.DB
}

func NewGroupRepo(db *gorm.DB) *GroupGorm {
	return &GroupGorm{
		db: db,
	}
}

func (g *GroupGorm) Fetch(ctx context.Context, req repository.GetGroupRequest) (ret []*entity.Group, err error) {
	logger := logrus.WithField("method", "repo.GroupGorm.Get")

	q := g.db.WithContext(ctx)

	if req.UserID != nil {
		q = q.Where("id = ?", req.UserID)
	}

	if req.PhoneNumber != nil {
		q = q.Where("phone_number = ?", req.PhoneNumber)
	}

	if req.Pin != nil {
		q = q.Where("pin = ?", req.Pin)
	}

	logger.Debugf("about to get data user")
	var item []Group
	if err = q.Find(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	group := Group{}
	return group.ToModels(item), nil
}

func (g *GroupGorm) Get(ctx context.Context, req repository.GetGroupRequest) (ret *entity.Group, err error) {
	logger := logrus.WithField("method", "repo.GroupGorm.Get")

	q := g.db.WithContext(ctx)

	if req.UserID != nil {
		q = q.Where("id = ?", req.UserID)
	}

	if req.PhoneNumber != nil {
		q = q.Where("phone_number = ?", req.PhoneNumber)
	}

	if req.Pin != nil {
		q = q.Where("pin = ?", req.Pin)
	}

	logger.Debugf("about to get data user")
	item := Group{}
	if err = q.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return item.ToModel(), nil
}

func (g *GroupGorm) Add(ctx context.Context, data entity.Group) (*entity.Group, error) {
	q := g.db.WithContext(ctx)

	group := Group{}
	item := group.FromModel(data)
	if err := q.Create(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return item.ToModel(), nil
}

func (g *GroupGorm) Update(ctx context.Context, data entity.Group) (ret *entity.Group, err error) {
	logger := logrus.WithField("method", "repo.UserGorm.UpdateUser")

	tx := g.db.WithContext(ctx)

	var group Group
	dataUpdate := group.FromModel(data)

	logger.Debugf("about to update user data")
	err = tx.Model(User{
		ID: data.ID,
	}).Updates(dataUpdate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		return nil, err
	}

	logger.Debugf("about to get user data")
	ret, err = g.Get(ctx, repository.GetGroupRequest{
		UserID: data.ID,
	})
	if err != nil {
		return nil, err
	}

	return
}

func (g *GroupGorm) Delete(ctx context.Context, id string) (bool, error) {
	q := g.db.WithContext(ctx)

	if err := q.Delete(&Group{}, id).Error; err != nil {
		return false, err
	}

	return true, nil
}

type Group struct {
	ID       *string `json:"id"`
	GroupKey *string `json:"group_key"`
	Name     *string `json:"name"`
	Level    *int    `json:"level"`
	Users    []User  `gorm:"foreignKey:GroupID;references:ID"`
}

func (g *Group) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == nil {
		g.ID = lo.ToPtr(uuid.New().String())

	}

	return
}

func (g *Group) FromModel(e entity.Group) *Group {
	return &Group{
		ID:       e.ID,
		GroupKey: e.GroupKey,
		Name:     e.Name,
		Level:    e.Level,
	}
}

func (g *Group) ToModel() *entity.Group {
	var users []entity.User
	for _, user := range g.Users {
		users = append(users, *user.ToModel())
	}
	m := &entity.Group{
		ID:       g.ID,
		GroupKey: g.GroupKey,
		Name:     g.Name,
		Level:    g.Level,
		Users:    users,
	}

	return m
}

func (g *Group) ToModels(d []Group) []*entity.Group {
	var m []*entity.Group
	for _, v := range d {
		m = append(m, v.ToModel())
	}
	return m
}
