package main

import (
	"context"
	"fmt"
	"github.com/delyke/gophermat_bonus_system/internal/app"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	a, err := app.New(appCtx)
	if err != nil {
		log.Println(err)
		return
	}
	_ = a.ShowConfig(appCtx)
}
