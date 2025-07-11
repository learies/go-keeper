package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/learies/go-keeper/internal/api/proto/auth/v1"
	"github.com/learies/go-keeper/internal/server/repository"
)

type AuthService struct {
	authv1.UnimplementedAuthServiceServer
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Ping — пример метода gRPC.
func (s *AuthService) Ping(ctx context.Context, req *authv1.PingRequest) (*authv1.PingResponse, error) {
	return &authv1.PingResponse{
		Response: "Pong: " + req.Message,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	// Создаём пользователя
	userID, err := s.userRepo.CreateUser(ctx, req.Email, string(hashedPassword))
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, "email already registered")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	return &authv1.RegisterResponse{UserId: userID}, nil
}

func (s *AuthService) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	// Получаем пользователя
	userID, hashedPassword, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	// Генерируем токен (заглушка)
	token := fmt.Sprintf("generated-jwt-token-for-%s", userID)

	return &authv1.LoginResponse{Token: token}, nil
}
