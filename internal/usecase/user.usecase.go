package usecase

import (
	"context"
	"time"

	"github.com/ramadhia/mnc-test/internal/model"

	"github.com/go-playground/validator"
	"github.com/shortlyst-ai/go-helper"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type UserUsecase interface {
	RegisterUser(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, error)
	LoginUser(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error)
	UpdateUser(ctx context.Context, user model.Claim, req UpdateProfileRequest) (*UpdateProfileResponse, error)
	GetUser(ctx context.Context, user model.Claim) (*model.User, error)
}

type RegisterUserRequest struct {
	FirstName   *string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    *string `json:"last_name" validate:"required,min=2,max=50"`
	PhoneNumber *string `json:"phone_number" validate:"required,min=2,max=13,numeric"`
	Address     *string `json:"address" validate:"omitempty,max=500"`
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
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
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
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,max=20"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,max=20"`
	Address   *string `json:"address,omitempty" validate:"omitempty,max=500"`
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
	UserID      *string    `json:"user_id"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	PhoneNumber *string    `json:"phone_number"`
	Address     *string    `json:"address"`
	UpdatedDate *time.Time `json:"updated_date"`
}

type LoginUserResponse struct {
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}
