package service

import (
	"google.golang.org/grpc"

	authv1 "github.com/learies/go-keeper/internal/api/proto/auth/v1"
)

func RegisterServices(server *grpc.Server) {
	authService := NewAuthService()
	authv1.RegisterAuthServiceServer(server, authService)
}
