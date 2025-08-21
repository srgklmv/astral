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

func (r repository) IsAdminTokenValid(token string) (bool, error) {
	panic("not implemented")
}

func (r repository) CreateUser(login, hashedPassword string, isAdmin bool) (user.User, error) {
	panic("not implemented")
}

func (r repository) GetByLogin(login string) (user.User, error) {
	panic("not implemented")
}

func (r repository) ValidatePassword(userID int, hashedPassword string) (bool, error) {
	panic("not implemented")
}

func (r repository) SaveAuthToken(userID int, token string) error {
	panic("not implemented")
}

func (r repository) DeleteToken(login string) (bool, error) {
	panic("not implemented")
}
