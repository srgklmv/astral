package repository

import (
	"context"
	"log/slog"

	"github.com/srgklmv/astral/pkg/logger"
)

func (r repository) DeleteAllUserTokens(ctx context.Context, login string) error {
	err := r.conn.QueryRowContext(
		ctx,
		`delete from auth_token where user_login = $1;`,
		login,
	).Err()
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r repository) DeleteToken(ctx context.Context, token string) error {
	err := r.conn.QueryRowContext(
		ctx,
		`delete from auth_token where token = $1 returning true;`,
		token,
	).Err()
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r repository) SaveAuthToken(ctx context.Context, login, token string) error {
	err := r.conn.QueryRowContext(
		ctx,
		`insert into auth_token(token, user_login) values ($1, $2)`,
		token,
		login,
	).Err()
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return err
	}

	return nil
}
