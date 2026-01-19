package bonus

import (
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
	orderService "github.com/delyke/gophermat_bonus_system/internal/service/bonus/order"
	userService "github.com/delyke/gophermat_bonus_system/internal/service/bonus/user"
	withdrawalService "github.com/delyke/gophermat_bonus_system/internal/service/bonus/withdrawal"
	"github.com/delyke/gophermat_bonus_system/internal/service/workers"
)

var _ def.BonusService = (*service)(nil)

type service struct {
	bonusRepository   repository.BonusRepository
	tokenIssuer       def.TokenIssuer
	userService       def.UserService
	orderService      def.OrderService
	withdrawalService def.WithdrawalService
	accrualWorker     workers.AccrualWorker
	logger            *logger.Logger
}

func NewService(
	bonusRepository repository.BonusRepository,
	tokenIssuer def.TokenIssuer,
	accrualWorker workers.AccrualWorker,
	logger *logger.Logger,
) *service {
	return &service{
		bonusRepository: bonusRepository,
		tokenIssuer:     tokenIssuer,
		accrualWorker:   accrualWorker,
		logger:          logger,
	}
}

// Users - сервис управления пользователями
func (s *service) Users() def.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.bonusRepository, s.tokenIssuer, s.logger)
	}
	return s.userService
}

// Orders - сервис управления заказами
func (s *service) Orders() def.OrderService {
	if s.orderService == nil {
		s.orderService = orderService.NewService(s.bonusRepository, s.accrualWorker, s.logger)
	}
	return s.orderService
}

// Withdrawals - сервис управления выводами пользователя
func (s *service) Withdrawals() def.WithdrawalService {
	if s.withdrawalService == nil {
		s.withdrawalService = withdrawalService.NewService(s.bonusRepository, s.logger)
	}
	return s.withdrawalService
}
