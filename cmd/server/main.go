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
	application, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("could not create application",
			slog.String("error", err.Error()),
			slog.Any("config", cfg))
		os.Exit(1)
	}

	// Запускаем приложение
	if err := application.Run(); err != nil {
		slog.Error("Server stopped with error", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Server stopped gracefully")
}
