package repository

import (
	"context"
	"log/slog"

	"github.com/srgklmv/astral/internal/domain/user"
	"github.com/srgklmv/astral/pkg/logger"
)

func (r repository) IsLoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool

	err := r.conn.QueryRowContext(ctx, `select true from users where login = $1`, login).Scan(&exists)
	if err != nil {
		logger.Error("QueryRowContext", slog.String("error", err.Error()))
	}

	return exists, nil
}

func (r repository) IsAdminTokenValid(ctx context.Context, token string) (bool, error) {
	panic("not implemented")
}

func (r repository) CreateUser(ctx context.Context, login, hashedPassword string, isAdmin bool) (user.User, error) {
	panic("not implemented")
}

func (r repository) GetByLogin(ctx context.Context, login string) (user.User, error) {
	panic("not implemented")
}

func (r repository) ValidatePassword(ctx context.Context, userID int, hashedPassword string) (bool, error) {
	panic("not implemented")
}

func (r repository) SaveAuthToken(ctx context.Context, userID int, token string) error {
	panic("not implemented")
}

func (r repository) DeleteToken(ctx context.Context, login string) (bool, error) {
	panic("not implemented")
}
