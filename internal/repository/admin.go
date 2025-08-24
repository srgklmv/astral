package repository

import (
	"context"
	"log/slog"

	"github.com/srgklmv/astral/pkg/logger"
)

func (r repository) IsAdminTokenValid(ctx context.Context, token string) (bool, error) {
	var t string

	err := r.conn.QueryRowContext(ctx, `select value from secrets where name = 'admin_token';`).Scan(&t)
	if err != nil {
		logger.Error("QueryRowContext", slog.String("error", err.Error()))
		return false, err
	}

	return t == token, nil
}
