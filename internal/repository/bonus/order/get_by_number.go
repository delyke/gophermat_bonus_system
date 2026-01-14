package order

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func (repo *repository) GetByNumber(ctx context.Context, number string) (*model.Order, error) {
	builderSelectOne := sq.Select(
		"uuid",
		"order_id",
		"order_status",
		"accrual",
		"user_uuid",
		"uploaded_at",
	).
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_id": number}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		order       repoModel.Order
		gUuid       uuid.UUID
		gOrderId    string
		gStatus     repoModel.OrderStatus
		gAccrual    *float64
		gUserUUID   uuid.UUID
		gUploadedAt time.Time
	)

	err = repo.pool.QueryRow(ctx, query, args...).Scan(&gUuid, &gOrderId, &gStatus, &gAccrual, &gUserUUID, &gUploadedAt)
	if err != nil {
		return nil, err
	}
	order.UUID = gUuid
	order.OrderID = gOrderId
	order.Status = gStatus
	order.Accrual = gAccrual
	order.UserUUID = gUserUUID
	order.UploadedAt = gUploadedAt

	return lo.ToPtr(converter.RepoOrderToService(order)), nil
}
