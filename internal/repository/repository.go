package repository

import (
	"database/sql"
)

type repository struct {
	conn *sql.DB
}

func New(conn *sql.DB) *repository {
	return &repository{conn: conn}
}
