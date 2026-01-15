package order

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func (repo *repository) ListPending(ctx context.Context, limit uint64) ([]*model.Order, error) {
	builderSelect := sq.Select(
		"uuid",
		"order_id",
		"order_status",
		"accrual",
		"user_uuid",
		"uploaded_at",
	).
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_status": []repoModel.OrderStatus{
			repoModel.OrderNew,
			repoModel.OrderProcessing,
		}}).
		Limit(limit)
	query, args, err := builderSelect.ToSql()
	if err != nil {
		logger.Error(ctx, "Ошибка при сборке sql запроса", zap.Error(err))
		return nil, err
	}

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "Ошибка при выборе списка заказов", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var (
		orders       []*repoModel.Order
		gID          uuid.UUID
		gOrderID     string
		gOrderStatus repoModel.OrderStatus
		gAccrual     *float64
		gUserUUID    uuid.UUID
		gUploadedAt  time.Time
	)

	for rows.Next() {
		err = rows.Scan(&gID, &gOrderID, &gOrderStatus, &gAccrual, &gUserUUID, &gUploadedAt)
		if err != nil {
			logger.Error(ctx, "Ошибка при переборке полученных заказов", zap.Error(err))
			return nil, err
		}
		orders = append(orders, &repoModel.Order{
			UUID:       gID,
			OrderID:    gOrderID,
			Status:     gOrderStatus,
			Accrual:    gAccrual,
			UserUUID:   gUserUUID,
			UploadedAt: gUploadedAt,
		})
	}

	return converter.RepoOrderListToService(orders), nil
}
