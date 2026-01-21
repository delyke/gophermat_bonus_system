package withdrawal

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (repo *repository) SumAmountByUserUUID(ctx context.Context, userUUID uuid.UUID) (float64, error) {
	sumBuilder := sq.Select("COALESCE(sum(amount), 0)").
		From("withdrawals").
		Where(sq.Eq{"user_uuid": userUUID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sumBuilder.ToSql()
	if err != nil {
		return 0, err
	}
	var amount float64
	err = repo.pool.QueryRow(ctx, query, args...).Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
