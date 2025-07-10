package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/learies/go-keeper/internal/config"
)

// App представляет основное приложение, содержащее gRPC сервер и конфигурацию
type App struct {
	grpcServer *grpc.Server
	cfg        *config.Config
}

// NewApp создает новый экземпляр приложения
func NewApp(cfg *config.Config) (*App, error) {
	grpcServer := grpc.NewServer()

	return &App{
		grpcServer: grpcServer,
		cfg:        cfg,
	}, nil
}

// Run запускает приложение и gRPC сервер
func (a *App) Run() error {
	// Формируем адрес для прослушивания из конфигурации
	addr := net.JoinHostPort(a.cfg.GRPC.Host, a.cfg.GRPC.Port)

	// Создаем listener для TCP соединений
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Создаем контекст, который будет отменен при получении SIGINT или SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запускаем горутину для graceful shutdown
	go func() {
		<-ctx.Done()
		slog.Info("Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
	}()

	// Логируем информацию о запуске сервера
	slog.Info("Starting gRPC server", "address", addr)

	// Запускаем gRPC сервер
	if err := a.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server failed: %v", err)
	}

	return nil
}
