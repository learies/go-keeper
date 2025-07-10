package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/learies/go-keeper/internal/client/app"
	"github.com/learies/go-keeper/internal/config"
)

func main() {
	// Загружаем конфигурацию приложения
	cfg := config.MustLoadConfig()

	// Настройка логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.ParseLogLevel(cfg.Log.Level),
	}))
	slog.SetDefault(logger)

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("could not create application",
			slog.String("error", err.Error()),
			slog.Any("config", cfg))
		os.Exit(1)
	}

	// Создаем новое приложение, передавая загруженную конфигурацию
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Запускаем приложение
	if err := application.Run(ctx); err != nil {
		slog.Error("Client stopped with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Client stopped gracefully")
}
