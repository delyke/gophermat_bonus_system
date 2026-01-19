package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (repo *repository) SetBalanceByUUID(ctx context.Context, userID uuid.UUID, amount float64) error {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("balance", amount).
		Where(sq.Eq{"uuid": userID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при построении запроса на обновление баланса пользователя", zap.Error(err))
		return err
	}

	res, err := repo.pool.Exec(ctx, query, args...)
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при обновлении баланса", zap.Error(err))
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
