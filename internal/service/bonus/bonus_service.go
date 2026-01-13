package bonus

import (
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
	userService "github.com/delyke/gophermat_bonus_system/internal/service/bonus/user"
)

var _ def.BonusService = (*service)(nil)

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

func (s *service) Users() def.UserService {
	return userService.NewService(s.bonusRepository, s.tokenIssuer)
}
