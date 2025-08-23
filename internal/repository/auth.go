package repository

import (
	"context"
	"log/slog"

	"github.com/srgklmv/astral/pkg/logger"
)

func (r repository) IsAuthTokenExists(ctx context.Context, token string) (bool, error) {
	var exists bool

	err := r.conn.QueryRowContext(
		ctx,
		`select coalesce((select true from auth_token where token = $1), false);`,
		token,
	).Scan(&exists)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return false, err
	}

	return exists, nil
}

func (r repository) GetUserLoginByAuthToken(ctx context.Context, token string) (string, error) {
	var login string

	err := r.conn.QueryRowContext(
		ctx,
		`select coalesce((select user_login from auth_token where token = $1), '');`,
		token,
	).Scan(&login)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return "", err
	}

	return login, nil
}
