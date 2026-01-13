package main

import (
	"context"
	"fmt"
	"github.com/delyke/gophermat_bonus_system/internal/app"
	"github.com/delyke/gophermat_bonus_system/internal/closer"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		log.Println(err)
		return
	}
	_ = a.ShowConfig(appCtx)

	err = a.Run(appCtx)
	if err != nil {
		log.Println(err)
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
