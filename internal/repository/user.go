package repository

import (
	"context"
	"log/slog"

	"github.com/srgklmv/astral/internal/domain/user"
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

func (r repository) IsAdminTokenValid(ctx context.Context, token string) (bool, error) {
	var t string

	err := r.conn.QueryRowContext(ctx, `select value from secrets where name = 'admin_token'`).Scan(&t)
	if err != nil {
		logger.Error("QueryRowContext", slog.String("error", err.Error()))
		return false, err
	}

	return t == token, nil
}

func (r repository) CreateUser(ctx context.Context, login, hashedPassword string, isAdmin bool) (user.User, error) {
	var user user.User

	err := r.conn.QueryRowContext(
		ctx,
		`insert into "user"(login, password, is_admin) values ($1, $2, $3) returning id, login, is_admin`,
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

func (r repository) GetUserByLogin(ctx context.Context, login string) (user.User, error) {
	var user user.User

	err := r.conn.QueryRowContext(
		ctx,
		`select id, login, is_admin from "user" where login = $1;`,
		login,
	).Scan(&user.ID, &user.Login, &user.IsAdmin)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return user, err
	}

	return user, nil
}

func (r repository) ValidatePassword(ctx context.Context, userID int, hashedPassword string) (bool, error) {
	var isValid bool

	err := r.conn.QueryRowContext(
		ctx,
		`select true from "user" where id = $1 and password = $2;`,
		userID,
		hashedPassword,
	).Scan(&isValid)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return false, err
	}

	return isValid, nil
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

func (r repository) DeleteToken(ctx context.Context, token string) (bool, error) {
	var deleted bool

	err := r.conn.QueryRowContext(
		ctx,
		`delete from auth_token where token = $1 returning true`,
		token,
	).Scan(&deleted)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return false, err
	}

	return deleted, nil
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
