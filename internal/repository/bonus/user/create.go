package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/postgres"
	"github.com/delyke/gophermat_bonus_system/internal/repository/converter"
	repoModel "github.com/delyke/gophermat_bonus_system/internal/repository/model"
)

func (repo *repository) Create(ctx context.Context, login, password string) (*model.User, error) {
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("login", "password").
		Values(login, password).
		Suffix("RETURNING uuid, login, password, balance")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, err
	}

	var (
		cUuid     uuid.UUID
		cLogin    string
		cPassword string
		cBalance  float64
	)

	err = repo.pool.QueryRow(ctx, query, args...).Scan(&cUuid, &cLogin, &cPassword, &cBalance)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return nil, model.ErrLoginTaken
		}
		logger.Error(ctx, "[REGISTRATION] Failed to create user:", zap.Error(err))
		return nil, err
	}

	insertedUser := &repoModel.User{
		UUID:     cUuid,
		Login:    cLogin,
		Password: cPassword,
		Balance:  cBalance,
	}

	return converter.RepoUserToService(insertedUser), nil
}
