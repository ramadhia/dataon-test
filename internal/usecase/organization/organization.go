package organization

import (
	"context"

	"github.com/ramadhia/dataon-test/internal/config"
	"github.com/ramadhia/dataon-test/internal/entity"
	"github.com/ramadhia/dataon-test/internal/provider"
	"github.com/ramadhia/dataon-test/internal/repository"
	"github.com/ramadhia/dataon-test/internal/usecase"

	"github.com/sirupsen/logrus"
)

type OrganizationImplement struct {
	config           config.Config
	organizationRepo repository.OrganizationRepository
}

func NewOrganization(p *provider.Provider) *OrganizationImplement {
	return &OrganizationImplement{
		config:           p.Config(),
		organizationRepo: p.OrganizationRepo(),
	}
}

func (t *OrganizationImplement) FetchOrganization(ctx context.Context, req usecase.FetchOrganizationRequest) ([]*entity.Organization, error) {
	logger := logrus.WithField("method", "usecase.TransactionImplement.FetchTransaction")

	logger.Debugf("about to fetch the transactions")
	data, err := t.organizationRepo.FetchOrganization(ctx, repository.FetchTransactionFilter{
		UserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []*entity.Organization{}, nil
	}

	return data, nil
}

func (t *OrganizationImplement) FetchComplete(ctx context.Context) (interface{}, error) {
	logger := logrus.WithField("method", "usecase.TransactionImplement.FetchTransaction")

	logger.Debugf("about to fetch the complete transactions")
	fetch, err := t.organizationRepo.FetchComplete(ctx, repository.FetchTransactionFilter{})
	if err != nil {
		return nil, err
	}

	return fetch, nil
}
