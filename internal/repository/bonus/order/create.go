package order

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/postgres"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func (repo *repository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	repoOrder := converter.ServiceOrderToRepo(*order)
	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("order_id", "order_status", "accrual", "user_uuid", "uploaded_at").
		Values(repoOrder.OrderID, repoOrder.Status, repoOrder.Accrual, repoOrder.UserUUID, repoOrder.UploadedAt).
		Suffix("RETURNING uuid, order_id, order_status, accrual, user_uuid, uploaded_at")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		cUuid       uuid.UUID
		cOrderId    string
		cStatus     repoModel.OrderStatus
		cAccrual    *float64
		cUserUUID   uuid.UUID
		cUploadedAt time.Time
	)

	err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&cUuid,
		&cOrderId,
		&cStatus,
		&cAccrual,
		&cUserUUID,
		&cUploadedAt,
	)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return nil, model.ErrOrderIDAlreadyExists
		}
		logger.Error(ctx, "[ORDER CREATE] Failed to create order:", zap.Error(err))
		return nil, err
	}

	repoOrder.UUID = cUuid
	repoOrder.OrderID = cOrderId
	repoOrder.Status = cStatus
	repoOrder.Accrual = cAccrual
	repoOrder.UserUUID = cUserUUID
	repoOrder.UploadedAt = cUploadedAt

	return lo.ToPtr(converter.RepoOrderToService(repoOrder)), nil
}
