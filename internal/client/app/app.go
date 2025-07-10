package app

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/learies/go-keeper/internal/client/service"
	"github.com/learies/go-keeper/internal/config"
)

type App struct {
	authClient *service.AuthClient
	cfg        *config.Config
}

// New создает новый экземпляр приложения
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

// Run запускает приложение и gRPC клиент
func (a *App) Run() error {
	// Создаем контекст, который будет отменен при получении SIGINT или SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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
