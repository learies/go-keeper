package service

import (
	"context"

	authv1 "github.com/learies/go-keeper/internal/api/proto/auth/v1"
)

type AuthService struct {
	authv1.UnimplementedAuthServiceServer // Встраиваем "заглушку" из сгенерированного кода
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Ping — пример метода gRPC.
func (s *AuthService) Ping(ctx context.Context, req *authv1.PingRequest) (*authv1.PingResponse, error) {
	return &authv1.PingResponse{
		Response: "Pong: " + req.Message,
	}, nil
}
