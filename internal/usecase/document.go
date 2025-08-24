package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/srgklmv/astral/internal/domain/document"
	userDomain "github.com/srgklmv/astral/internal/domain/user"
	"github.com/srgklmv/astral/internal/models/apperrors"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
	"github.com/srgklmv/astral/pkg/utils"
)

func (u usecase) UploadDocument(ctx context.Context, token string, meta dto.UploadDocumentRequestMetadata, json dto.UploadDocumentRequestJSON, file *bytes.Buffer) (dto.APIResponse[any, *dto.UploadFileResponse], int) {
	isTokenValid, user, err := u.authorizeUserByToken(ctx, token)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.UploadFileResponse](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}
	if !isTokenValid {
		return dto.NewAPIResponse[any, *dto.UploadFileResponse](&dto.Error{
			Code: apperrors.UnauthorizedErrorCode,
			Text: apperrors.UnauthorizedErrorText,
		}, nil, nil), http.StatusUnauthorized
	}

	isMetaValid, errorText := document.ValidateDocumentMetadata(meta)
	if !isMetaValid {
		return dto.NewAPIResponse[any, *dto.UploadFileResponse](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: errorText,
		}, nil, nil), http.StatusBadRequest
	}

	if meta.IsFile && len(file.Bytes()) == 0 {
		return dto.NewAPIResponse[any, *dto.UploadFileResponse](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.FileNotProvidedErrorText,
		}, nil, nil), http.StatusBadRequest
	}

	doc, err := u.documentRepository.UploadDocument(
		ctx,
		user.Login,
		meta.Name,
		meta.IsFile,
		meta.Mimetype,
		meta.IsPublic,
		meta.GrantedTo,
		json,
		file,
	)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.UploadFileResponse](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	return dto.NewAPIResponse[any, *dto.UploadFileResponse](nil, nil, &dto.UploadFileResponse{
		JSON:     doc.JSON,
		Filename: doc.Filename,
	}), http.StatusCreated
}

func (u usecase) GetDocuments(ctx context.Context, token, login, filterKey, filterValue string, limit int) (dto.APIResponse[any, *dto.GetDocumentsResponse], int) {
	panic("not implemented")
}

func (u usecase) GetDocumentsHead(ctx context.Context, token, login, filterKey, filterValue string, limit int) (bool, int) {
	panic("not implemented")
}

func (u usecase) GetDocument(ctx context.Context, token, id string) (dto.APIResponse[any, any], []byte, map[string]string, int) {
	var isTokenValid bool
	var user userDomain.User
	var err error

	if token != "" {
		isTokenValid, user, err = u.authorizeUserByToken(ctx, token)
		if err != nil {
			logger.Error("repository call error", slog.String("error", err.Error()))
			return dto.NewAPIResponse[any, any](&dto.Error{
				Code: apperrors.InternalErrorErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil), nil, nil, http.StatusInternalServerError
		}
		if !isTokenValid {
			return dto.NewAPIResponse[any, any](&dto.Error{
				Code: apperrors.UnauthorizedErrorCode,
				Text: apperrors.UnauthorizedErrorText,
			}, nil, nil), nil, nil, http.StatusUnauthorized
		}
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.BadIDProvidedErrorText,
		}, nil, nil), nil, nil, http.StatusBadRequest
	}

	doc, err := u.documentRepository.GetDocument(ctx, uid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.DocumentNotFoundErrorText,
		}, nil, nil), nil, nil, http.StatusNotFound
	}
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), nil, nil, http.StatusInternalServerError
	}

	if !doc.IsPublic && token == "" {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.ForbiddenErrorCode,
			Text: apperrors.ForbiddenErrorText,
		}, nil, nil), nil, nil, http.StatusForbidden
	}

	isOwner := utils.IsSliceIncludesValue(doc.GrantedTo, user.Login)
	if !isOwner && doc.Owner != user.Login && !user.IsAdmin {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.ForbiddenErrorCode,
			Text: apperrors.ForbiddenErrorText,
		}, nil, nil), nil, nil, http.StatusForbidden
	}

	if !doc.IsFile {
		return dto.NewAPIResponse[any, any](
			nil,
			nil,
			&doc.JSON,
		), nil, nil, http.StatusOK
	}

	headers := make(map[string]string)
	headers["Content-Type"] = doc.Mimetype
	headers["Content-Disposition"] = fmt.Sprintf(`attachment; filename="%s"`, doc.Filename)
	headers["Content-Transfer-Encoding"] = "binary"
	headers["Content-Length"] = fmt.Sprintf("%d", len(doc.File))
	headers["Last-Modified"] = doc.CreatedAt.Format("Mon, 02 Jan 2006 15:04:05 GMT")

	return dto.APIResponse[any, any]{}, doc.File, headers, http.StatusOK
}

func (u usecase) GetDocumentHead(ctx context.Context, token, id string) (dto.APIResponse[any, any], []byte, map[string]string, int) {
	var isTokenValid bool
	var user userDomain.User
	var err error

	if token != "" {
		isTokenValid, user, err = u.authorizeUserByToken(ctx, token)
		if err != nil {
			logger.Error("repository call error", slog.String("error", err.Error()))
			return dto.NewAPIResponse[any, any](&dto.Error{
				Code: apperrors.InternalErrorErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil), nil, nil, http.StatusInternalServerError
		}
		if !isTokenValid {
			return dto.NewAPIResponse[any, any](&dto.Error{
				Code: apperrors.UnauthorizedErrorCode,
				Text: apperrors.UnauthorizedErrorText,
			}, nil, nil), nil, nil, http.StatusUnauthorized
		}
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.BadIDProvidedErrorText,
		}, nil, nil), nil, nil, http.StatusBadRequest
	}

	doc, err := u.documentRepository.GetDocument(ctx, uid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.DocumentNotFoundErrorText,
		}, nil, nil), nil, nil, http.StatusNotFound
	}
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), nil, nil, http.StatusInternalServerError
	}

	if !doc.IsPublic && token == "" {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.ForbiddenErrorCode,
			Text: apperrors.ForbiddenErrorText,
		}, nil, nil), nil, nil, http.StatusForbidden
	}

	isOwner := utils.IsSliceIncludesValue(doc.GrantedTo, user.Login)
	if !isOwner && doc.Owner != user.Login && !user.IsAdmin {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.ForbiddenErrorCode,
			Text: apperrors.ForbiddenErrorText,
		}, nil, nil), nil, nil, http.StatusForbidden
	}

	if !doc.IsFile {
		return dto.NewAPIResponse[any, any](
			nil,
			nil,
			&doc.JSON,
		), nil, nil, http.StatusOK
	}

	headers := make(map[string]string)
	headers["Content-Type"] = doc.Mimetype
	headers["Content-Disposition"] = fmt.Sprintf(`attachment; filename="%s"`, doc.Filename)
	headers["Content-Transfer-Encoding"] = "binary"
	headers["Content-Length"] = fmt.Sprintf("%d", len(doc.File))
	headers["Last-Modified"] = doc.CreatedAt.Format("Mon, 02 Jan 2006 15:04:05 GMT")

	return dto.APIResponse[any, any]{}, doc.File, headers, http.StatusOK
}

func (u usecase) DeleteDocument(ctx context.Context, token, id string) (dto.APIResponse[any, *dto.DeleteDocumentResponse], int) {
	if token == "" {
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.UnauthorizedErrorCode,
			Text: apperrors.UnauthorizedErrorText,
		}, nil, nil), http.StatusUnauthorized
	}
	if id == "" {
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.DocumentIDNotProvidedErrorText,
		}, nil, nil), http.StatusBadRequest
	}

	isValid, user, err := u.authorizeUserByToken(ctx, token)
	if err != nil {
		logger.Error("validateAuthToken error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}
	if !isValid {
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.UnauthorizedErrorCode,
			Text: apperrors.UnauthorizedErrorText,
		}, nil, nil), http.StatusUnauthorized
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Error("uuid parse error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.BadIDProvidedErrorText,
		}, nil, nil), http.StatusBadRequest
	}

	doc, err := u.documentRepository.GetDocument(ctx, uid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.DocumentNotFoundErrorText,
		}, nil, nil), http.StatusNotFound
	}
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	isUserGrantedToDocument := utils.IsSliceIncludesValue(doc.GrantedTo, id)

	if !isUserGrantedToDocument && doc.Owner != user.Login && !user.IsAdmin {
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.ForbiddenErrorCode,
			Text: apperrors.ForbiddenErrorText,
		}, nil, nil), http.StatusForbidden
	}

	err = u.documentRepository.DeleteDocument(ctx, doc.ID)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	return dto.NewAPIResponse[any, *dto.DeleteDocumentResponse](nil, dto.DeleteDocumentResponse{
		id: true,
	}, nil), http.StatusOK
}
