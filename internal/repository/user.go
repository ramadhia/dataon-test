package repository

import (
	"context"

	"github.com/ramadhia/dataon-test/internal/entity"
)

type UserRepository interface {
	Register(ctx context.Context, data entity.User) (*entity.User, error)
	FetchUser(ctx context.Context, req FetchUserRequest) ([]*entity.User, error)
	GetUser(ctx context.Context, req GetUserRequest) (*entity.User, error)
	UpdateUser(ctx context.Context, data entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
}

type GetUserRequest struct {
	UserID      *string
	Pin         *string
	PhoneNumber *string
}

type FetchUserRequest struct {
	OrganizationID *string
}
