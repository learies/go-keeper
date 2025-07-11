package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learies/go-keeper/internal/config"
	"github.com/learies/go-keeper/internal/server"
	"github.com/learies/go-keeper/internal/server/repository"
	"github.com/learies/go-keeper/internal/server/service"
)

// App представляет основное приложение, содержащее gRPC сервер и конфигурацию
type App struct {
	grpcServer *grpc.Server
	cfg        *config.Config
	pool       *pgxpool.Pool
}

// NewApp создает новый экземпляр приложения
func New(cfg *config.Config) (*App, error) {
	// Подключаемся к PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.Database.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	userRepo := repository.NewUserRepository(pool)
	grpcServer := server.NewGRPCServer()
	service.RegisterServices(grpcServer, userRepo)

	return &App{
		grpcServer: grpcServer,
		cfg:        cfg,
		pool:       pool,
	}, nil
}

// Run запускает приложение и gRPC сервер
func (a *App) Run() error {
	// Формируем адрес для прослушивания из конфигурации
	addr := net.JoinHostPort(a.cfg.GRPC.Host, a.cfg.GRPC.Port)

	// Создаем listener для TCP соединений
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Создаем контекст, который будет отменен при получении SIGINT или SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Запускаем горутину для graceful shutdown
	go func() {
		<-ctx.Done()
		slog.Info("Shutting down gRPC server gracefully...")
		a.Shutdown()
	}()

	// Логируем информацию о запуске сервера
	slog.Info("Starting gRPC server", slog.String("address", addr))

	// Запускаем gRPC сервер
	if err := a.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server failed: %v", err)
	}

	return nil
}

func (a *App) Shutdown() {
	a.grpcServer.GracefulStop()
	a.pool.Close()
}
