package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (repo *repository) GetBalanceByUUID(ctx context.Context, uuidUser uuid.UUID) (float64, error) {
	builderSelectOne := sq.Select(
		"balance",
	).
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"uuid": uuidUser}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return 0, err
	}

	var gBalance float64

	err = repo.pool.QueryRow(ctx, query, args...).Scan(&gBalance)
	if err != nil {
		return 0, err
	}

	return gBalance, nil
}
