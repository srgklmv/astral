package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/srgklmv/astral/internal/domain/document"
	"github.com/srgklmv/astral/internal/models/apperrors"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
	"github.com/srgklmv/astral/pkg/utils"
)

func (u usecase) UploadDocument(ctx context.Context, token string, meta dto.UploadDocumentRequestMetadata, json dto.UploadDocumentRequestJSON, file *bytes.Buffer) (dto.APIResponse[any, *dto.UploadFileResponse], int) {
	isTokenValid, login, err := u.validateAuthToken(ctx, token)
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

	doc, err := u.documentRepository.UploadDocument(
		ctx,
		login,
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

func (u usecase) GetDocument(ctx context.Context, token, id string) (dto.APIResponse[any, *dto.GetDocumentResponse], int) {
	panic("not implemented")
}

func (u usecase) GetDocumentHead(ctx context.Context, token, id string) (bool, int) {
	panic("not implemented")
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

	isValid, login, err := u.validateAuthToken(ctx, token)
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

	if !isUserGrantedToDocument && doc.Owner != login {
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
