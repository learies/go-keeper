package server

import (
	"google.golang.org/grpc"

	"github.com/learies/go-keeper/internal/server/interceptors"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor,
			interceptors.RecoveryInterceptor,
		),
	)
}
