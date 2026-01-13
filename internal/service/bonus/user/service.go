package user

import (
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	bonusRepository repository.BonusRepository
	tokenIssuer     def.TokenIssuer
}

func NewService(bonusRepository repository.BonusRepository, tokenIssuer def.TokenIssuer) *service {
	return &service{
		bonusRepository: bonusRepository,
		tokenIssuer:     tokenIssuer,
	}
}
