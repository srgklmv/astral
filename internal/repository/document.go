package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
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
	var id uuid.UUID

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
	).Scan(&id, &jsonb, &doc.Filename)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return document.Document{}, err
	}

	for _, grantedToLogin := range grantedTo {
		err = tx.QueryRowContext(
			ctx,
			`insert into user_document_access (user_login, document_id) values ($1, $2);`,
			grantedToLogin,
			id,
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

	return doc, nil
}

func (r repository) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.ExecContext(
		ctx,
		`delete from document where id=$1;`,
		id.String(),
	)
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r repository) GetDocument(ctx context.Context, id uuid.UUID) (document.Document, error) {
	var doc document.Document
	var uid string
	var jsonb, owners []byte

	err := r.conn.QueryRowContext(
		ctx,
		`select d.id, d.name, d.is_file, d.is_public, d.mimetype, d.json, d.file, d.created_at, d.owner_login, to_json(array_agg(uda.user_login)) as owners
		from document d
		left join public.user_document_access uda on d.id = uda.document_id
		where d.id = $1
		group by d.id;`,
		id.String(),
	).Scan(&uid, &doc.Filename, &doc.IsFile, &doc.IsPublic, &doc.Mimetype, &jsonb, &doc.File, &doc.CreatedAt, &doc.Owner, &owners)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return doc, err
	}
	if err != nil {
		logger.Error("QueryRowContext error", slog.String("error", err.Error()))
		return doc, err
	}

	err = json.Unmarshal(jsonb, &doc.JSON)
	if err != nil {
		logger.Error("Unmarshal error", slog.String("error", err.Error()))
		return doc, err
	}

	err = json.Unmarshal(owners, &doc.GrantedTo)
	if err != nil {
		logger.Error("Unmarshal error", slog.String("error", err.Error()))
		return doc, err
	}

	doc.ID, err = uuid.Parse(uid)
	if err != nil {
		logger.Error("uuid.FromBytes error", slog.String("error", err.Error()))
		return doc, err
	}

	return doc, nil
}

func (r repository) GetDocumentsData(ctx context.Context, userLogin string, isAdmin bool, login, key, value string, limit int) (document.DocumentsData, error) {
	query := []string{
		`select d.id, d.name, d.is_file, d.is_public, d.mimetype, d.created_at, to_json(array_agg(uda.user_login)) as owners
		from document d
		left join user_document_access uda on d.id = uda.document_id`,
		`where`,
	}
	var args []any

	switch {
	case login == "":
		query = append(query, `(d.owner_login = $1 or uda.user_login = $1)`)
		args = append(args, userLogin)
	case isAdmin:
		query = append(query, `(d.owner_login = $1 or uda.user_login = $1)`)
		args = append(args, login)
	default:
		query = append(query, `(d.owner_login = $1 or uda.user_login = $1) and 
			(d.is_public or d.owner_login = $2 or uda.user_login = $2)`)
		args = append(args, login, userLogin)
	}

	if key != "" {
		query = append(query, `and`)
		query = append(query, fmt.Sprintf(`(d.%s = $%d)`, key, len(args)+1))
		args = append(args, value)
	}

	query = append(
		query,
		`group by d.id, d.name, d.created_at
		order by d.name ASC, d.created_at ASC`,
	)
	query = append(query, fmt.Sprintf("limit $%d", len(args)+1))
	args = append(args, limit)

	q := strings.Join(query, " ")
	var docs document.DocumentsData

	rows, err := r.conn.QueryContext(ctx, q, args...)
	if err != nil {
		logger.Error("QueryContext error", slog.String("error", err.Error()))
		return docs, err
	}
	defer rows.Close()

	for rows.Next() {
		var doc document.Data
		var buf []byte

		err = rows.Scan(
			&doc.ID,
			&doc.Filename,
			&doc.IsFile,
			&doc.IsPublic,
			&doc.Mimetype,
			&doc.CreatedAt,
			&buf,
		)
		if err != nil {
			logger.Error("QueryContext error", slog.String("error", err.Error()))
			return docs, err
		}

		err = json.Unmarshal(buf, &doc.GrantedTo)
		if err != nil {
			logger.Error("Unmarshal error", slog.String("error", err.Error()))
			return docs, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}
