package bonus

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
)

func (r *repository) CreateWithdrawalWithBalance(ctx context.Context, withdrawal *model.Withdrawal) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			return
		}
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			r.logger.Error(ctx, "Ошибка при откате транзакции списания", zap.Error(rollbackErr))
		}
	}()

	var balance float64
	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE uuid = $1 FOR UPDATE", withdrawal.UserUUID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrUserNotFound
		}
		return err
	}

	if balance < withdrawal.Amount {
		return model.ErrNotEnoughBalance
	}

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

	err = tx.QueryRow(ctx, query, args...).Scan(&withdrawal.UUID)
	if err != nil {
		r.logger.Error(ctx, "Ошибка при вставке нового вывода", zap.Error(err))
		return err
	}

	newBalance := balance - withdrawal.Amount
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("balance", newBalance).
		Where(sq.Eq{"uuid": withdrawal.UserUUID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		r.logger.Error(ctx, "Ошибка при построении запроса на обновление баланса пользователя", zap.Error(err))
		return err
	}

	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error(ctx, "Ошибка при обновлении баланса", zap.Error(err))
		return err
	}

	if res.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
