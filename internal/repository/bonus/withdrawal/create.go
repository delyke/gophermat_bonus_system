package withdrawal

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
)

func (repo *repository) Create(ctx context.Context, withdrawal *model.Withdrawal) error {
	rWithdrawal := converter.ServiceWithdrawalToRepo(*withdrawal)
	builderInsert := sq.Insert("withdrawals").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "order_id", "amount", "processed_at").
		Values(rWithdrawal.UserUUID, rWithdrawal.OrderID, rWithdrawal.Amount, rWithdrawal.ProcessedAt).
		Suffix("RETURNING uuid")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	err = repo.pool.QueryRow(ctx, query, args...).Scan(&withdrawal.UUID)
	if err != nil {
		repo.logger.Error(ctx, "Ошибка при вставке нового вывода")
		return err
	}

	return nil
}
