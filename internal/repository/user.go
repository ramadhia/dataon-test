package repository

import (
	"context"

	"github.com/ramadhia/mnc-test/internal/model"
)

type UserRepository interface {
	Register(ctx context.Context, data model.User) (*model.User, error)
	GetUser(ctx context.Context, req GetUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, data model.User) (*model.User, error)
}

type GetUserRequest struct {
	UserID      *string
	Pin         *string
	PhoneNumber *string
}
