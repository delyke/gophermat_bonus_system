package withdrawal

import (
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
)

var _ def.WithdrawalService = (*service)(nil)

type service struct {
	bonusRepository repository.BonusRepository
}

func NewService(bonusRepository repository.BonusRepository) *service {
	return &service{
		bonusRepository: bonusRepository,
	}
}
