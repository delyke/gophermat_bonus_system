package bonus

import (
	repository2 "github.com/delyke/gophermat_bonus_system/internal/repository"
	usersRepo "github.com/delyke/gophermat_bonus_system/internal/repository/bonus/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository2.BonusRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Users() repository2.UserRepository {
	return usersRepo.NewRepository(r.pool)
}
