package app

import (
	"context"
	"fmt"
	"github.com/delyke/gophermat_bonus_system/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"

	security "github.com/delyke/gophermat_bonus_system/internal/api"
	bonusV1Api "github.com/delyke/gophermat_bonus_system/internal/api/bonus/v1"
	"github.com/delyke/gophermat_bonus_system/internal/closer"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	bonusRepository "github.com/delyke/gophermat_bonus_system/internal/repository/bonus"
	"github.com/delyke/gophermat_bonus_system/internal/service"
	bonusService "github.com/delyke/gophermat_bonus_system/internal/service/bonus"
	tokenIssuer "github.com/delyke/gophermat_bonus_system/internal/service/bonus/user/jwt"
	"github.com/delyke/gophermat_bonus_system/internal/service/workers"
	accrualV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/accrual/v1"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

type diContainer struct {
	bonusV1API      bonusV1.Handler
	bonusService    service.BonusService
	pool            *pgxpool.Pool
	poolCfg         *pgxpool.Config
	bonusServer     *bonusV1.Server
	bonusRepository repository.BonusRepository
	sec             *security.SecurityHandler
	tokenIssuer     service.TokenIssuer
	accrualClient   *accrualV1.Client
	accrualWorker   workers.AccrualWorker
	logger          *logger.Logger
}

func NewDIContainer(appLogger *logger.Logger) *diContainer {

	return &diContainer{logger: appLogger}
}

func (di *diContainer) AccrualWorker(ctx context.Context) workers.AccrualWorker {
	if di.accrualWorker == nil {
		p := workers.NewAccrualProcessor(
			ctx,
			di.logger,
			di.BonusRepository(ctx),
			di.AccrualClient(ctx), config.Get().Accrual.WorkersCount(),
		)

		closer.AddNamed("Accrual Worker", func(ctx context.Context) error {
			p.Stop()
			return nil
		})

		di.accrualWorker = p
	}
	return di.accrualWorker
}

func (di *diContainer) AccrualClient(_ context.Context) *accrualV1.Client {
	if di.accrualClient == nil {
		client, err := accrualV1.NewClient(config.Get().Accrual.SystemAddress())
		if err != nil {
			return nil
		}
		di.accrualClient = client
	}
	return di.accrualClient
}

func (di *diContainer) TokenIssuer() service.TokenIssuer {
	if di.tokenIssuer == nil {
		di.tokenIssuer = tokenIssuer.New(config.Get().JWT.Secret())
	}
	return di.tokenIssuer
}

func (di *diContainer) Sec() *security.SecurityHandler {
	if di.sec == nil {
		di.sec = security.New(config.Get().JWT.Secret())
	}
	return di.sec
}

func (di *diContainer) BonusServer(ctx context.Context) *bonusV1.Server {
	if di.bonusServer == nil {
		bonusServer, err := bonusV1.NewServer(di.BonusV1Api(ctx), di.Sec())
		if err != nil {
			panic(fmt.Sprintf("Error creating bonus server: %v", err))
		}
		di.bonusServer = bonusServer
	}
	return di.bonusServer
}

func (di *diContainer) BonusV1Api(ctx context.Context) bonusV1.Handler {
	if di.bonusV1API == nil {
		di.bonusV1API = bonusV1Api.NewAPI(di.BonusService(ctx), di.logger)
	}
	return di.bonusV1API
}

func (di *diContainer) BonusService(ctx context.Context) service.BonusService {
	if di.bonusService == nil {
		di.bonusService = bonusService.NewService(di.BonusRepository(ctx), di.TokenIssuer(), di.AccrualWorker(ctx), di.logger)
	}
	return di.bonusService
}

func (di *diContainer) BonusRepository(ctx context.Context) repository.BonusRepository {
	if di.bonusRepository == nil {
		di.bonusRepository = bonusRepository.NewRepository(di.PostgresPool(ctx), di.logger)
	}
	return di.bonusRepository
}

func (di *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if di.pool == nil {
		pool, err := pgxpool.New(ctx, config.Get().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres: %w", err))
		}
		closer.AddNamed("Postgres Pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})
		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to ping postgres: %w", err))
		}
		di.pool = pool
	}
	return di.pool
}

func (di *diContainer) PoolCfg(_ context.Context) *pgxpool.Config {
	if di.poolCfg == nil {
		poolCfg, err := pgxpool.ParseConfig(config.Get().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres: %w", err))
		}
		di.poolCfg = poolCfg
	}
	return di.poolCfg
}
