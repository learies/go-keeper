package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создаёт нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, email, hashedPassword string) (string, error) {
	var userID string
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		email, hashedPassword,
	).Scan(&userID)
	if err != nil {
		slog.Error("Database operation failed", "error", err)
		if isDuplicateKeyError(err) {
			return "", ErrUserExists
		}
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	slog.Debug("User created", "user_id", userID)
	return userID, nil
}

// GetUserByEmail возвращает хеш пароля пользователя
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (id string, hashedPassword string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT id, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&id, &hashedPassword)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrUserNotFound
	}

	return id, hashedPassword, err
}

func isDuplicateKeyError(err error) bool {
	const pgErrCodeUniqueViolation = "23505"
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation
}
