package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/app"
	"github.com/delyke/gophermat_bonus_system/internal/closer"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/logger"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(fmt.Errorf("error loading config: %w", err))
	}
	appLogger, err := logger.New(config.Get().Logger.Level(), config.Get().Logger.AsJSON())
	if err != nil {
		panic(fmt.Errorf("error initializing logger: %w", err))
	}
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown(appLogger)

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx, appLogger)
	if err != nil {
		appLogger.Fatal(appCtx, "error starting application", zap.Error(err))
		return
	}
	err = a.ShowConfig(appCtx)
	if err != nil {
		appLogger.Fatal(appCtx, "error showing config", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		appLogger.Fatal(appCtx, "error running app", zap.Error(err))
		return
	}
}

func gracefulShutdown(appLogger *logger.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Get().App.ShutdownTimeout())
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		appLogger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
