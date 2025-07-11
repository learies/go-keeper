package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/learies/go-keeper/internal/api/proto/auth/v1"
	"github.com/learies/go-keeper/internal/server/repository"
)

func RegisterServices(server *grpc.Server, userRepo *repository.UserRepository) {
	authService := NewAuthService(userRepo)
	authv1.RegisterAuthServiceServer(server, authService)

	reflection.Register(server)
}
