package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/closer"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/migrator"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	router      *chi.Mux
	logger      *logger.Logger
}

func New(ctx context.Context, appLogger *logger.Logger) (*App, error) {
	a := &App{logger: appLogger}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initMigrator,
		a.initCloser,
		a.initWorkers,
		a.initRouter,
		a.initServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer(a.logger)
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(a.logger)
	return nil
}

func (a *App) initMigrator(ctx context.Context) error {
	poolCfg, err := pgxpool.ParseConfig(config.Get().Postgres.URI())
	if err != nil {
		return err
	}

	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*poolCfg.ConnConfig), config.Get().Postgres.MigrationDirectory(), a.logger)
	err = migratorRunner.Up(ctx)
	if err != nil {
		return err
	}
	a.logger.Info(ctx, "migrations applied")
	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Mount("/", a.diContainer.BonusServer(ctx))
	a.router = r
	return nil
}

func (a *App) initServer(_ context.Context) error {
	server := &http.Server{
		Addr:              config.Get().HTTP.RunAddress(),
		Handler:           a.router,
		ReadHeaderTimeout: config.Get().HTTP.ReadTimeout(),
	}
	a.httpServer = server
	return nil
}

func (a *App) initWorkers(ctx context.Context) error {
	p := a.diContainer.AccrualWorker(ctx)

	if err := p.Bootstrap(ctx); err != nil {
		return err
	}
	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	a.logger.Info(ctx, "starting http server", zap.String("RunAddress:", config.Get().HTTP.RunAddress()))
	closer.AddNamed("HTTP Server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) ShowConfig(ctx context.Context) error {
	a.logger.Debug(ctx, "App configuration:", zap.Any("config", config.Get()))
	a.logger.Debug(ctx,
		"Logger config:",
		zap.Any("asJson", config.Get().Logger.AsJSON()),
		zap.Any("level", config.Get().Logger.Level()),
	)

	a.logger.Debug(ctx,
		"Accrual config:",
		zap.Any("SystemAddress", config.Get().Accrual.SystemAddress()),
	)

	a.logger.Debug(ctx,
		"Bonus HTTP Config:",
		zap.String("RunAddress:", config.Get().HTTP.RunAddress()),
		zap.Any("ReadTimeout:", config.Get().HTTP.ReadTimeout()),
	)

	a.logger.Debug(ctx,
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

	a.logger.Debug(ctx,
		"JWT Config:",
		zap.Any("Secret", config.Get().JWT.Secret()))
	return nil
}
