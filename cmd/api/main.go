package main

import (
	"context"
	"github.com/novabankapp/wallet.api/infrastructure/di"
	"github.com/novabankapp/wallet.api/server"
	"github.com/swaggo/swag/example/basic/docs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	docs.SwaggerInfo.Title = "Novabank Wallet API"
	docs.SwaggerInfo.Description = "This is an API for Managing Novabank Wallets."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "api/v1"

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	container := di.BuildContainer()
	if err := container.Invoke(func(app server.App) {
		app.Run(ctx)

	}); err != nil {

		log.Fatalf("main %s", err.Error())
	}

	<-ctx.Done()
}
