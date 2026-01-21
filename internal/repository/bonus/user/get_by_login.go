package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func (repo *repository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	builderSelectOne := sq.Select(
		"uuid",
		"login",
		"password",
		"balance",
	).
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"login": login}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		user      repoModel.User
		gUUID     uuid.UUID
		gLogin    string
		gPassword string
		gBalance  float64
	)

	err = repo.pool.QueryRow(ctx, query, args...).Scan(&gUUID, &gLogin, &gPassword, &gBalance)
	if err != nil {
		return nil, err
	}

	user.UUID = gUUID
	user.Login = login
	user.Password = gPassword
	user.Balance = gBalance

	return converter.RepoUserToService(&user), nil
}
