package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

// LoggingInterceptor — для логирования запросов.
func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	// Вызываем обработчик
	resp, err := handler(ctx, req)

	// Логируем информацию о запросе
	slog.Info(
		"gRPC request",
		"method", info.FullMethod,
		"duration", time.Since(start),
		"error", err,
	)

	return resp, err
}

// RecoveryInterceptor — перехватывает паники
func RecoveryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("recovered from panic", slog.Any("panic", r))
			err = fmt.Errorf("internal server error")
		}
	}()

	resp, err = handler(ctx, req)
	return
}
