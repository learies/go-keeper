package service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/learies/go-keeper/internal/api/proto/auth/v1"
	"github.com/learies/go-keeper/internal/config"
)

type AuthClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(cfg *config.Config) (*AuthClient, error) {
	addr := cfg.GRPC.Host + ":" + cfg.GRPC.Port

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthServiceClient(conn)

	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AuthClient) Ping(ctx context.Context, message string) (string, error) {
	resp, err := c.client.Ping(ctx, &authv1.PingRequest{Message: message})
	if err != nil {
		return "", err
	}

	return resp.Response, nil
}

func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
