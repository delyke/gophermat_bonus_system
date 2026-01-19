package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
)

func (repo *repository) SetStatusByUUID(ctx context.Context, uuid uuid.UUID, status model.OrderStatus) error {
	newStatus := converter.ServiceStatusToRepo(status)

	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("order_status", newStatus).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при построении запроса на обновление статуса заказа", zap.Error(err))
		return err
	}

	res, err := repo.pool.Exec(ctx, query, args...)
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при обновлении статуса заказа", zap.Error(err))
		return err
	}
	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}
	return nil
}
