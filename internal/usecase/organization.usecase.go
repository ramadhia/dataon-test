package usecase

import (
	"context"

	"github.com/ramadhia/dataon-test/internal/entity"
)

type OrganizationUsecase interface {
	FetchOrganization(ctx context.Context, args FetchOrganizationRequest) ([]*entity.Organization, error)
	FetchComplete(ctx context.Context) (interface{}, error)
}

type FetchOrganizationRequest struct {
	UserID *string
}

type FetchOrganizationResponse struct {
	Result interface{} `json:"result"`
}
