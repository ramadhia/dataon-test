package usecase

import (
	"context"
	"time"

	"github.com/ramadhia/dataon-test/internal/entity"

	"github.com/go-playground/validator"
	"github.com/shortlyst-ai/go-helper"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type UserUsecase interface {
	RegisterUser(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, error)
	// LoginUser(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error)
	UpdateUser(ctx context.Context, entity entity.User) (*UpdateProfileResponse, error)
	GetUser(ctx context.Context, entity entity.User) (*entity.User, error)
	FetchUser(ctx context.Context, req FetchUserRequest) ([]*entity.User, error)
}

type RegisterUserRequest struct {
	GroupID     *string `json:"group_id" validate:"required,min=2,max=50"`
	EmployeeID  *string `json:"employee_id" validate:"required,min=2,max=50"`
	Name        *string `json:"name" validate:"required,min=2,max=50"`
	PhoneNumber *string `json:"phone_number" validate:"required,min=2,max=13,numeric"`
	Pin         *string `json:"pin" validate:"required,len=6,numeric"`
}

func (r *RegisterUserRequest) Validate() (err error) {
	validate = validator.New()
	if err := validate.Struct(r); err != nil {
		errMsg := err.Error()
		return helper.NewParameterError(&errMsg)
	}
	return
}

type RegisterUserResponse struct {
	UserID      *string    `json:"user_id"`
	GroupID     *string    `json:"group_id"`
	EmployeeID  *string    `json:"employee_id"`
	Name        *string    `json:"Name"`
	PhoneNumber *string    `json:"phone_number"`
	Address     *string    `json:"address"`
	CreatedDate *time.Time `json:"created_date"`
}

type LoginUserRequest struct {
	PhoneNumber *string `json:"phone_number" validate:"required,min=2,max=13,numeric"`
	Pin         *string `json:"pin" validate:"required,min=1,numeric"`
}

func (l *LoginUserRequest) Validate() (err error) {
	validate = validator.New()
	if err := validate.Struct(l); err != nil {
		errMsg := err.Error()
		return helper.NewParameterError(&errMsg)
	}
	return
}

type UpdateProfileRequest struct {
	ID          *string `json:"id" validate:"required"`
	GroupID     *string `json:"group_id" validate:"required,min=2,max=50"`
	EmployeeID  *string `json:"employee_id" validate:"required,min=2,max=50"`
	Name        *string `json:"first_name,omitempty" validate:"omitempty,max=20"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,max=500"`
}

func (r *UpdateProfileRequest) Validate() (err error) {
	validate = validator.New()
	if err := validate.Struct(r); err != nil {
		errMsg := err.Error()
		return helper.NewParameterError(&errMsg)
	}
	return
}

type UpdateProfileResponse struct {
	UserID      *string `json:"user_id"`
	GroupID     *string `json:"group_id" validate:"required,min=2,max=50"`
	EmployeeID  *string `json:"employee_id" validate:"required,min=2,max=50"`
	Name        *string `json:"first_name,omitempty" validate:"omitempty,max=20"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,max=500"`
}

type LoginUserResponse struct {
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}

type FetchUserRequest struct {
	OrganizationID *string `json:"group_id" validate:"required,min=2,max=50"`
}
