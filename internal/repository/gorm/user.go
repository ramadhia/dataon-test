package gorm

import (
	"context"
	"errors"
	"github.com/ramadhia/mnc-test/internal/repository"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/ramadhia/mnc-test/internal/model"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UserGorm struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserGorm {
	return &UserGorm{
		db: db,
	}
}

func (u *UserGorm) GetUser(ctx context.Context, req repository.GetUserRequest) (ret *model.User, err error) {
	logger := logrus.WithField("method", "repo.UserGorm.RegisterUser")

	q := u.db.WithContext(ctx)

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
	item := User{}
	if err = q.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return item.ToModel(), nil
}

func (u *UserGorm) Register(ctx context.Context, data model.User) (*model.User, error) {
	q := u.db.WithContext(ctx)

	user := User{}
	item := user.FromModel(data)
	if err := q.Create(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return item.ToModel(), nil
}

func (u *UserGorm) UpdateUser(ctx context.Context, data model.User) (ret *model.User, err error) {
	logger := logrus.WithField("method", "repo.UserGorm.UpdateUser")

	tx := u.db.WithContext(ctx)

	var user User
	dataUpdate := user.FromModel(data)

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
	ret, err = u.GetUser(ctx, repository.GetUserRequest{
		UserID: data.ID,
	})
	if err != nil {
		return nil, err
	}

	return
}

type User struct {
	ID          *string    `json:"id"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	PhoneNumber *string    `json:"phone_number"`
	Address     *string    `json:"address"`
	Pin         *string    `json:"pin"`
	CreatedDate *time.Time `json:"created_date"`
	UpdatedDate *time.Time `json:"updated_date"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == nil {
		u.ID = lo.ToPtr(uuid.New().String())

	}

	if u.CreatedDate == nil {
		u.CreatedDate = lo.ToPtr(time.Now())
	}

	return
}

func (u *User) FromModel(model model.User) *User {
	return &User{
		ID:          model.ID,
		FirstName:   model.FirstName,
		LastName:    model.LastName,
		PhoneNumber: model.PhoneNumber,
		Pin:         model.Pin,
		Address:     model.Address,
		CreatedDate: model.CreatedDate,
		UpdatedDate: model.UpdatedDate,
	}
}

func (u *User) ToModel() *model.User {
	m := &model.User{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address,
		Pin:         u.Pin,
		CreatedDate: u.CreatedDate,
		UpdatedDate: u.UpdatedDate,
	}

	return m
}

func (u *User) ToModels(d []User) []*model.User {
	var m []*model.User
	for _, v := range d {
		m = append(m, v.ToModel())
	}
	return m
}
