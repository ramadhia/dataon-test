package gorm

import (
	"context"
	"errors"
	"time"

	"github.com/ramadhia/dataon-test/internal/entity"
	"github.com/ramadhia/dataon-test/internal/repository"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
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

func (u *UserGorm) GetUser(ctx context.Context, req repository.GetUserRequest) (ret *entity.User, err error) {
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

func (u *UserGorm) FetchUser(ctx context.Context, req repository.FetchUserRequest) ([]*entity.User, error) {
	logger := logrus.WithField("method", "repo.UserGorm.FetchUser")

	q := u.db.WithContext(ctx)

	if req.OrganizationID != nil {
		q = q.Where("organization_id = ?", req.OrganizationID)
	}

	logger.Debugf("about to get data user")
	var item []User
	if err := q.Find(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	user := User{}
	return user.ToModels(item), nil
}

func (u *UserGorm) Register(ctx context.Context, data entity.User) (*entity.User, error) {
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

func (u *UserGorm) UpdateUser(ctx context.Context, data entity.User) (ret *entity.User, err error) {
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

func (u *UserGorm) DeleteUser(ctx context.Context, id string) (bool, error) {
	logger := logrus.WithField("method", "repo.UserGorm.DeleteUser")

	q := u.db.WithContext(ctx)

	logger.Debugf("about to delete user data")
	if err := q.Delete(&Group{}, id).Error; err != nil {
		return false, err
	}

	return true, nil
}

type User struct {
	ID          *string    `json:"id"`
	GroupID     *string    `json:"group_id"`
	EmployeeID  *string    `json:"employee_id"`
	Name        *string    `json:"Name"`
	PhoneNumber *string    `json:"phone_number"`
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

func (u *User) FromModel(model entity.User) *User {
	return &User{
		ID:          model.ID,
		GroupID:     model.GroupID,
		EmployeeID:  model.EmployeeID,
		Name:        model.Name,
		PhoneNumber: model.PhoneNumber,
		Pin:         model.Pin,
		CreatedDate: model.CreatedDate,
	}
}

func (u *User) ToModel() *entity.User {
	m := &entity.User{
		ID:          u.ID,
		GroupID:     u.GroupID,
		EmployeeID:  u.EmployeeID,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Pin:         u.Pin,
		CreatedDate: u.CreatedDate,
	}

	return m
}

func (u *User) ToModels(d []User) []*entity.User {
	var m []*entity.User
	for _, v := range d {
		m = append(m, v.ToModel())
	}
	return m
}
