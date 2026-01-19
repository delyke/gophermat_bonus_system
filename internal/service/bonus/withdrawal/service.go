package withdrawal

import (
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
)

var _ def.WithdrawalService = (*service)(nil)

type service struct {
	bonusRepository repository.BonusRepository
	logger          *logger.Logger
}

func NewService(bonusRepository repository.BonusRepository, logger *logger.Logger) *service {
	return &service{
		bonusRepository: bonusRepository,
		logger:          logger,
	}
}
