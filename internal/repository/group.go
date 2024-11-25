package repository

import (
	"context"

	"github.com/ramadhia/dataon-test/internal/entity"
)

type GroupRepository interface {
	Add(ctx context.Context, data entity.Group) (*entity.Group, error)
	Get(ctx context.Context, req GetGroupRequest) (*entity.Group, error)
	Update(ctx context.Context, data entity.Group) (*entity.Group, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type GetGroupRequest struct {
	UserID      *string
	Pin         *string
	PhoneNumber *string
}
