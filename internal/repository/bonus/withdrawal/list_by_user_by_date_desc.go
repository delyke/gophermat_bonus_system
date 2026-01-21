package withdrawal

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func (repo *repository) ListByUserByDateDesc(ctx context.Context, userUUID uuid.UUID) ([]*model.Withdrawal, error) {
	builderSelect := sq.Select(
		"uuid",
		"user_uuid",
		"order_id",
		"amount",
		"processed_at",
	).
		From("withdrawals").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_uuid": userUUID}).
		OrderBy("processed_at DESC")

	query, args, err := builderSelect.ToSql()
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при построении запроса на получение списка выводов", zap.Error(err))
		return nil, err
	}

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при выборке списка выводов", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var (
		withdrawals  []*repoModel.Withdrawal
		gID          uuid.UUID
		gUserUUID    uuid.UUID
		gOrderID     string
		gAmount      float64
		gProcessedAt time.Time
	)

	for rows.Next() {
		err = rows.Scan(&gID, &gUserUUID, &gOrderID, &gAmount, &gProcessedAt)
		if err != nil {
			repo.logger.Error(ctx, "Ошибка при сканировании выданных данных в Repo Withdrawal", zap.Error(err))
			return nil, err
		}
		withdrawals = append(withdrawals, &repoModel.Withdrawal{
			UUID:        gID,
			UserUUID:    gUserUUID,
			OrderID:     gOrderID,
			Amount:      gAmount,
			ProcessedAt: gProcessedAt,
		})
	}
	return converter.RepoWithdrawalListToService(withdrawals), nil
}
