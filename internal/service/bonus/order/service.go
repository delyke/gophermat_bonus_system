package order

import (
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
	"github.com/delyke/gophermat_bonus_system/internal/service/workers"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	bonusRepository repository.BonusRepository
	accrualWorker   workers.AccrualWorker
	logger          *logger.Logger
}

func NewService(
	bonusRepository repository.BonusRepository,
	accrualWorker workers.AccrualWorker,
	logger *logger.Logger,
) *service {
	return &service{
		bonusRepository: bonusRepository,
		accrualWorker:   accrualWorker,
		logger:          logger,
	}
}
