package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/srgklmv/astral/pkg/logger"
)

func New(host, port, database, user, password string) (*sql.DB, error) {
	data := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		database,
	)

	conn, err := sql.Open("postgres", data)
	if err != nil {
		logger.Error("connection to database failed", slog.String("error", err.Error()))
		return nil, err
	}

	err = conn.QueryRowContext(context.Background(), "select 1;").Err()
	if err != nil {
		logger.Error("db test query failed", slog.String("error", err.Error()))
		return nil, err
	}

	return conn, nil
}

func Migrate(conn *sql.DB, path string, version int) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		logger.Error("migrations driver set up error", slog.String("error", err.Error()))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"postgres", driver)
	if err != nil {
		logger.Error("migrate instance creation error", slog.String("error", err.Error()))
		return err
	}

	err = m.Migrate(uint(version))
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("migrations up error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func Shutdown(conn *sql.DB) error {
	return conn.Close()
}

func SeedAdminToken(conn *sql.DB, adminToken string) error {
	if adminToken == "" {
		adminToken = "test"
	}

	err := conn.QueryRowContext(
		context.Background(),
		"insert into secrets (name, value) values ('admin_token', $1) on conflict (name) do update set value = $1;",
		adminToken,
	).Err()
	if err != nil {
		logger.Error("db admin seeding failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}
