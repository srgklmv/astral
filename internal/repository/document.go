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

	err = r.conn.QueryRowContext(
		ctx,
		`insert into document (name, is_file, is_public, mimetype, json, file, owner_login) values ($1, $2, $3, $4, $5, $6, $7) returning json, name;`,
		filename,
		isFile,
		isPublic,
		mimetype,
		jsonb,
		file.Bytes(),
		login,
	).Scan(&jsonb, &doc.Filename)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return document.Document{}, err
	}

	err = json.Unmarshal(jsonb, &jsonM)
	if err != nil {
		logger.Error("Unmarshal error", slog.String("error", err.Error()))
		return document.Document{}, err
	}

	doc.JSON = jsonM

	// TODO: Grants to.

	return doc, nil
}
