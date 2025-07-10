package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/learies/go-keeper/internal/client/service"
	"github.com/learies/go-keeper/internal/config"
)

type App struct {
	authClient *service.AuthClient
	cfg        *config.Config
}

func New(cfg *config.Config) (*App, error) {
	// Создаем gRPC клиент для сервиса аутентификации
	authClient, err := service.NewAuthClient(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:        cfg,
		authClient: authClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, err := a.authClient.Ping(ctx, "Hello!")
	if err != nil {
		slog.Error("ping failed",
			slog.String("error", err.Error()))
		return err
	}

	slog.Info("ping successful",
		slog.String("response", response))
	return nil
}

func (a *App) Close() error {
	if a.authClient != nil {
		return a.authClient.Close()
	}
	return nil
}
