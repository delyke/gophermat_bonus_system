package app

import (
	"context"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/logger"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.Get().Logger.Level(),
		config.Get().Logger.AsJson(),
	)
}

func (a *App) ShowConfig(ctx context.Context) error {
	logger.Debug(ctx, "App configuration:", zap.Any("config", config.Get()))
	logger.Debug(ctx,
		"Logger config:",
		zap.Any("asJson", config.Get().Logger.AsJson()),
		zap.Any("level", config.Get().Logger.Level()),
	)

	logger.Debug(ctx,
		"Accrual config:",
		zap.Any("SystemAddress", config.Get().Accrual.SystemAddress()),
	)

	logger.Debug(ctx,
		"Bonus HTTP Config:",
		zap.String("RunAddress:", config.Get().HTTP.RunAddress()),
		zap.Any("ReadTimeout:", config.Get().HTTP.ReadTimeout()),
	)

	logger.Debug(ctx,
		"PostgresConfig:",
		zap.String("host", config.Get().Postgres.Host()),
		zap.Int("port", config.Get().Postgres.Port()),
		zap.String("user", config.Get().Postgres.User()),
		zap.String("password", config.Get().Postgres.Password()),
		zap.String("database", config.Get().Postgres.Database()),
		zap.String("Database URI", config.Get().Postgres.DatabaseURI()),
		zap.String("URI", config.Get().Postgres.URI()),
		zap.String("Migrations dir", config.Get().Postgres.MigrationDirectory()),
	)
	return nil
}
