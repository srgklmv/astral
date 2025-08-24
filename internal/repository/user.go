package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	userDomain "github.com/srgklmv/astral/internal/domain/user"
	"github.com/srgklmv/astral/pkg/logger"
)

func (r repository) IsLoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool

	err := r.conn.QueryRowContext(ctx, `select coalesce((select true from "user" where login = $1), false);`, login).Scan(&exists)
	if err != nil {
		logger.Error("QueryRowContext", slog.String("error", err.Error()))
		return false, err
	}

	return exists, nil
}

func (r repository) CreateUser(ctx context.Context, login, hashedPassword string, isAdmin bool) (userDomain.User, error) {
	var user userDomain.User

	err := r.conn.QueryRowContext(
		ctx,
		`insert into "user"(login, password, is_admin) values ($1, $2, $3) returning id, login, is_admin;`,
		login,
		hashedPassword,
		isAdmin,
	).Scan(&user.ID, &user.Login, &user.IsAdmin)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return user, err
	}

	return user, nil
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

func (r repository) GetUserHashedPassword(ctx context.Context, login string) (string, error) {
	var password string

	err := r.conn.QueryRowContext(
		ctx,
		`select password from "user" where login = $1;`,
		login,
	).Scan(&password)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return "", err
	}

	return password, nil
}
