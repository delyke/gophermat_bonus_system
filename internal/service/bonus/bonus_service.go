package bonus

import (
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	def "github.com/delyke/gophermat_bonus_system/internal/service"
	orderService "github.com/delyke/gophermat_bonus_system/internal/service/bonus/order"
	userService "github.com/delyke/gophermat_bonus_system/internal/service/bonus/user"
	withdrawalService "github.com/delyke/gophermat_bonus_system/internal/service/bonus/withdrawal"
)

var _ def.BonusService = (*service)(nil)

type service struct {
	bonusRepository   repository.BonusRepository
	tokenIssuer       def.TokenIssuer
	userService       def.UserService
	orderService      def.OrderService
	withdrawalService def.WithdrawalService
}

func NewService(bonusRepository repository.BonusRepository, tokenIssuer def.TokenIssuer) *service {
	return &service{
		bonusRepository: bonusRepository,
		tokenIssuer:     tokenIssuer,
	}
}

// Users - сервис управления пользователями
func (s *service) Users() def.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.bonusRepository, s.tokenIssuer)
	}
	return s.userService
}

// Orders - сервис управления заказами
func (s *service) Orders() def.OrderService {
	if s.orderService == nil {
		s.orderService = orderService.NewService(s.bonusRepository)
	}
	return s.orderService
}

// Withdrawals - сервис управления выводами пользователя
func (s *service) Withdrawals() def.WithdrawalService {
	if s.withdrawalService == nil {
		s.withdrawalService = withdrawalService.NewService(s.bonusRepository)
	}
	return s.withdrawalService
}
