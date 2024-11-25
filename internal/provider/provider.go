package provider

import (
	"github.com/ramadhia/dataon-test/internal/config"
	"github.com/ramadhia/dataon-test/internal/repository"
	"github.com/ramadhia/dataon-test/internal/service"
	"github.com/ramadhia/dataon-test/internal/usecase"

	"gorm.io/gorm"
)

type Provider struct {
	config config.Config
	db     *gorm.DB

	// initialize service
	service ServiceProvider

	// initialize usecase
	usecase UsecaseProvider

	// initialize repository
	repo RepositoryProvider
}

type ServiceProvider struct {
	messageBus service.MessageBus
}

type UsecaseProvider struct {
	user         usecase.UserUsecase
	organization usecase.OrganizationUsecase
}

type RepositoryProvider struct {
	user         repository.UserRepository
	group        repository.GroupRepository
	organization repository.OrganizationRepository
}

func NewProvider() *Provider {

	return &Provider{}
}

func (p *Provider) Config() config.Config {
	return p.config
}

func (p *Provider) SetConfig(c config.Config) {
	p.config = c
}

func (p *Provider) DB() *gorm.DB {
	return p.db
}

func (p *Provider) SetDB(d *gorm.DB) {
	p.db = d
}

func (p *Provider) MessageBus() service.MessageBus {
	return p.service.messageBus
}

func (p *Provider) SetMessageBus(m service.MessageBus) {
	p.service.messageBus = m
}

func (p *Provider) UserUseCase() usecase.UserUsecase {
	return p.usecase.user
}

func (p *Provider) SetUserUseCase(u usecase.UserUsecase) {
	p.usecase.user = u
}

func (p *Provider) OrganizationUseCase() usecase.OrganizationUsecase {
	return p.usecase.organization
}

func (p *Provider) SetOrganizationUseCase(u usecase.OrganizationUsecase) {
	p.usecase.organization = u
}

func (p *Provider) UserRepo() repository.UserRepository {
	return p.repo.user
}

func (p *Provider) SetUserRepo(r repository.UserRepository) {
	p.repo.user = r
}
func (p *Provider) GroupRepo() repository.GroupRepository {
	return p.repo.group
}

func (p *Provider) SetGroupRepo(r repository.GroupRepository) {
	p.repo.group = r
}

func (p *Provider) OrganizationRepo() repository.OrganizationRepository {
	return p.repo.organization
}

func (p *Provider) SetOrganizationRepo(r repository.OrganizationRepository) {
	p.repo.organization = r
}
