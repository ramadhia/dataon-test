package provider

import (
	"github.com/ramadhia/mnc-test/internal/config"
	"github.com/ramadhia/mnc-test/internal/repository"
	"github.com/ramadhia/mnc-test/internal/service"
	"github.com/ramadhia/mnc-test/internal/usecase"

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
	user        usecase.UserUsecase
	transaction usecase.TransactionUsecase
	algo        usecase.AlgoUsecase
}

type RepositoryProvider struct {
	user        repository.UserRepository
	transaction repository.TransactionRepository
	balance     repository.BalanceRepository
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

func (p *Provider) TransactionUseCase() usecase.TransactionUsecase {
	return p.usecase.transaction
}

func (p *Provider) SetTransactionUseCase(u usecase.TransactionUsecase) {
	p.usecase.transaction = u
}

func (p *Provider) AlgoUseCase() usecase.AlgoUsecase {
	return p.usecase.algo
}

func (p *Provider) SetAlgoUseCase(u usecase.AlgoUsecase) {
	p.usecase.algo = u
}

func (p *Provider) UserRepo() repository.UserRepository {
	return p.repo.user
}

func (p *Provider) SetUserRepo(r repository.UserRepository) {
	p.repo.user = r
}

func (p *Provider) TransactionRepo() repository.TransactionRepository {
	return p.repo.transaction
}

func (p *Provider) SetTransactionRepo(r repository.TransactionRepository) {
	p.repo.transaction = r
}

func (p *Provider) BalanceRepo() repository.BalanceRepository {
	return p.repo.balance
}

func (p *Provider) SetBalanceRepo(r repository.BalanceRepository) {
	p.repo.balance = r
}
