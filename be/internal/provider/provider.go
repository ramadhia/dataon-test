package provider

import (
	"github.com/ramadhia/bosnet/be/internal/config"
	"github.com/ramadhia/bosnet/be/internal/repository"
	"github.com/ramadhia/bosnet/be/internal/usecase"

	"gorm.io/gorm"
)

type Provider struct {
	config config.Config
	db     *gorm.DB

	// initialize usecase
	historyUc usecase.HistoryUsecase

	// initialize repository
	historyRepo repository.HistoryRepository
	balanceRepo repository.BalanceRepository
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

func (p *Provider) HistoryUseCase() usecase.HistoryUsecase {
	return p.historyUc
}

func (p *Provider) SetHistoryUseCase(u usecase.HistoryUsecase) {
	p.historyUc = u
}

func (p *Provider) HistoryRepo() repository.HistoryRepository {
	return p.historyRepo
}

func (p *Provider) SetHistoryRepo(r repository.HistoryRepository) {
	p.historyRepo = r
}

func (p *Provider) BalanceRepo() repository.BalanceRepository {
	return p.balanceRepo
}

func (p *Provider) SetBalanceRepo(r repository.BalanceRepository) {
	p.balanceRepo = r
}
