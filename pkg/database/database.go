package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

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

	db, err := sql.Open("postgres", data)
	if err != nil {
		logger.Error("connection to database failed", slog.String("error", err.Error()))
		return nil, err
	}

	err = db.QueryRowContext(context.Background(), "select 1;").Err()
	if err != nil {
		logger.Error("db test query failed", slog.String("error", err.Error()))
		return nil, err
	}

	return db, nil
}
