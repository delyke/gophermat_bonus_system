package bonus

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
	def "github.com/delyke/gophermat_bonus_system/internal/repository"
	ordersRepo "github.com/delyke/gophermat_bonus_system/internal/repository/bonus/order"
	usersRepo "github.com/delyke/gophermat_bonus_system/internal/repository/bonus/user"
	withdrawalsRepo "github.com/delyke/gophermat_bonus_system/internal/repository/bonus/withdrawal"
)

var _ def.BonusRepository = (*repository)(nil)

type repository struct {
	pool           *pgxpool.Pool
	usersRepo      def.UserRepository
	orderRepo      def.OrderRepository
	withdrawalRepo def.WithdrawalRepository
	logger         *logger.Logger
}

func NewRepository(pool *pgxpool.Pool, appLogger *logger.Logger) *repository {
	return &repository{
		pool:   pool,
		logger: appLogger,
	}
}

func (r *repository) Users() def.UserRepository {
	if r.usersRepo == nil {
		r.usersRepo = usersRepo.NewRepository(r.pool, r.logger)
	}
	return r.usersRepo
}

func (r *repository) Orders() def.OrderRepository {
	if r.orderRepo == nil {
		r.orderRepo = ordersRepo.NewRepository(r.pool, r.logger)
	}
	return r.orderRepo
}

func (r *repository) Withdrawals() def.WithdrawalRepository {
	if r.withdrawalRepo == nil {
		r.withdrawalRepo = withdrawalsRepo.NewRepository(r.pool, r.logger)
	}
	return r.withdrawalRepo
}
