package order

import (
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/delyke/gophermat_bonus_system/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	pool   *pgxpool.Pool
	logger *logger.Logger
}

func NewRepository(pool *pgxpool.Pool, l *logger.Logger) *repository {
	return &repository{
		pool:   pool,
		logger: l,
	}
}
