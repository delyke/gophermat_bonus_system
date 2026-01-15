package bonus

import (
	"github.com/jackc/pgx/v5/pgxpool"

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
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Users() def.UserRepository {
	if r.usersRepo == nil {
		r.usersRepo = usersRepo.NewRepository(r.pool)
	}
	return r.usersRepo
}

func (r *repository) Orders() def.OrderRepository {
	if r.orderRepo == nil {
		r.orderRepo = ordersRepo.NewRepository(r.pool)
	}
	return r.orderRepo
}

func (r *repository) Withdrawals() def.WithdrawalRepository {
	if r.withdrawalRepo == nil {
		r.withdrawalRepo = withdrawalsRepo.NewRepository(r.pool)
	}
	return r.withdrawalRepo
}
