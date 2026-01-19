package user

import (
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	bonusRepository repository.BonusRepository
	tokenIssuer     def.TokenIssuer
	logger          *logger.Logger
}

func NewService(
	bonusRepository repository.BonusRepository,
	tokenIssuer def.TokenIssuer,
	logger *logger.Logger) *service {
	return &service{
		bonusRepository: bonusRepository,
		tokenIssuer:     tokenIssuer,
		logger:          logger,
	}
}
