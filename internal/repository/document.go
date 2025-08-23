package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/srgklmv/astral/internal/domain/document"
	"github.com/srgklmv/astral/pkg/logger"
)

func (r repository) UploadDocument(
	ctx context.Context,
	login, filename string,
	isFile bool,
	mimetype string,
	isPublic bool,
	grantedTo []string,
	jsonM map[string]any,
	file *bytes.Buffer,
) (document.Document, error) {
	var doc document.Document

	jsonb, err := json.Marshal(jsonM)
	if err != nil {
		logger.Error("json.Marshal error", slog.String("error", err.Error()))
		return document.Document{}, err
	}

	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("database transaction error", slog.String("error", err.Error()))
		return document.Document{}, err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
			if err != nil {
				logger.Error("commit error", slog.String("error", err.Error()))
			}
			return
		}

		err = tx.Rollback()
		if err != nil {
			logger.Error("rollback error", slog.String("error", err.Error()))
		}
	}()

	err = tx.QueryRowContext(
		ctx,
		`insert into document (name, is_file, is_public, mimetype, json, file, owner_login) values ($1, $2, $3, $4, $5, $6, $7) returning id, json, name;`,
		filename,
		isFile,
		isPublic,
		mimetype,
		jsonb,
		file.Bytes(),
		login,
	).Scan(&doc.ID, &jsonb, &doc.Filename)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return document.Document{}, err
	}

	for _, grantedToLogin := range grantedTo {
		err = tx.QueryRowContext(
			ctx,
			`insert into user_document_access (user_login, document_id) values ($1, $2);`,
			grantedToLogin,
			doc.ID,
		).Err()
		if err != nil {
			// TODO: Add validation for granted to logins.
			logger.Error("QueryRowContext error", slog.String("error", err.Error()))
			return document.Document{}, err
		}
	}

	err = json.Unmarshal(jsonb, &jsonM)
	if err != nil {
		logger.Error("Unmarshal error", slog.String("error", err.Error()))
		return document.Document{}, err
	}

	doc.JSON = jsonM
	doc.ID = 0

	return doc, nil
}
