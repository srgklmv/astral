package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	userDomain "github.com/srgklmv/astral/internal/domain/user"
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

func (r repository) GetUserByAuthToken(ctx context.Context, token string) (userDomain.User, error) {
	var user userDomain.User

	err := r.conn.QueryRowContext(
		ctx,
		`select u.id, u.login, u.is_admin 
		from auth_token at
		left join "user" u on at.user_login = u.login		    
		where at.token = $1;`,
		token,
	).Scan(&user.ID, &user.Login, &user.IsAdmin)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		}
		return user, err
	}

	return user, nil
}
