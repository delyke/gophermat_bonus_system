package order

import (
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
	"github.com/delyke/gophermat_bonus_system/internal/service/workers"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	bonusRepository repository.BonusRepository
	accrualWorker   workers.AccrualWorker
}

func NewService(bonusRepository repository.BonusRepository, accrualWorker workers.AccrualWorker) *service {
	return &service{
		bonusRepository: bonusRepository,
		accrualWorker:   accrualWorker,
	}
}
