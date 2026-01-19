package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (repo *repository) SetAccrualByUUID(ctx context.Context, uuid uuid.UUID, accrual float64) error {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("accrual", accrual).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при построении запроса на обновление начисленного вознаграждения", zap.Error(err))
		return err
	}

	res, err := repo.pool.Exec(ctx, query, args...)
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при записи суммы вознаграждения")
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}
	return nil
}
