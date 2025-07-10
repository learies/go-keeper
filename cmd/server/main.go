package main

import (
	"log/slog"
	"os"

	"github.com/learies/go-keeper/internal/config"
	"github.com/learies/go-keeper/internal/server/app"
)

// main - точка входа в приложение
func main() {
	// Загружаем конфигурацию приложения
	cfg := config.MustLoadConfig()

	// Настройка логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.ParseLogLevel(cfg.Log.Level),
	}))
	slog.SetDefault(logger)

	// Создаем новое приложение, передавая загруженную конфигурацию
	app, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("Could not create application", slog.Any("error", err))
		os.Exit(1)
	}

	// Запускаем приложение
	if err := app.Run(); err != nil {
		slog.Error("Server stopped with error", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Server stopped gracefully")
}
