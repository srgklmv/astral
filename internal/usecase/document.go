package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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

	isMetaValid, errorText := u.validateDocumentMetadata(meta)
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

	if !meta.IsFile && len(json) == 0 {
		return dto.NewAPIResponse[any, *dto.UploadFileResponse](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.JSONNotProvidedErrorText,
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

func (u usecase) GetDocuments(ctx context.Context, request dto.GetDocumentsRequest) (dto.APIResponse[any, *dto.GetDocumentsResponse], int) {
	isAuthorized, user, err := u.authorizeUserByToken(ctx, request.Token)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.GetDocumentsResponse](&dto.Error{
			Code: apperrors.AuthInternalErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}
	if !isAuthorized {
		return dto.NewAPIResponse[any, *dto.GetDocumentsResponse](&dto.Error{
			Code: apperrors.UnauthorizedErrorCode,
			Text: apperrors.UnauthorizedErrorText,
		}, nil, nil), http.StatusUnauthorized
	}

	request, isValid, errText := u.validateGetDocumentsRequest(request)
	if !isValid {
		return dto.NewAPIResponse[any, *dto.GetDocumentsResponse](&dto.Error{
			Code: apperrors.AuthInternalErrorCode,
			Text: errText,
		}, nil, nil), http.StatusInternalServerError
	}

	// Strict passing dto into repository is not kinda safe, but
	// $ symbol doing things.
	documents, err := u.documentRepository.GetDocumentsData(
		ctx,
		user.Login,
		user.IsAdmin,
		request.Login,
		request.Key,
		request.Value,
		request.Limit,
	)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, *dto.GetDocumentsResponse](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	docsDTO := dto.NewGetDocumentsResponse().FromDomain(documents)

	return dto.NewAPIResponse[any, *dto.GetDocumentsResponse](nil, nil, &docsDTO), http.StatusOK
}

func (u usecase) GetDocument(ctx context.Context, token, id string) (dto.APIResponse[any, any], []byte, string, int) {
	isAuthorized, user, err := u.authorizeUserByToken(ctx, token)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.AuthInternalErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), nil, "", http.StatusInternalServerError
	}
	if !isAuthorized {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.UnauthorizedErrorCode,
			Text: apperrors.UnauthorizedErrorText,
		}, nil, nil), nil, "", http.StatusUnauthorized
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.BadIDProvidedErrorText,
		}, nil, nil), nil, "", http.StatusBadRequest
	}

	doc, err := u.documentRepository.GetDocument(ctx, uid)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.DocumentNotFoundErrorText,
		}, nil, nil), nil, "", http.StatusNotFound
	}
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), nil, "", http.StatusInternalServerError
	}

	if !doc.IsPublic && token == "" {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.ForbiddenErrorCode,
			Text: apperrors.ForbiddenErrorText,
		}, nil, nil), nil, "", http.StatusForbidden
	}

	isOwner := utils.IsSliceIncludesValue(doc.GrantedTo, user.Login)
	if !isOwner && doc.Owner != user.Login && !user.IsAdmin {
		return dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.ForbiddenErrorCode,
			Text: apperrors.ForbiddenErrorText,
		}, nil, nil), nil, "", http.StatusForbidden
	}

	if !doc.IsFile {
		return dto.NewAPIResponse[any, any](
			nil,
			nil,
			&doc.JSON,
		), nil, "", http.StatusOK
	}

	return dto.APIResponse[any, any]{}, doc.File, doc.Filename, http.StatusOK
}

func (u usecase) DeleteDocument(ctx context.Context, token, id string) (dto.APIResponse[any, *dto.DeleteDocumentResponse], int) {
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

func (u usecase) validateDocumentMetadata(metadata dto.UploadDocumentRequestMetadata) (bool, apperrors.ErrorText) {
	// May be moved to config.
	allowedMimeTypes := map[string]bool{
		"application/pdf": true,
		"image/png":       true,
		"image/jpeg":      true,
		"image/webp":      true,
		"image/gif":       true,
	}

	if !allowedMimeTypes[metadata.Mimetype] {
		return false, apperrors.MimeTypeNotAllowedErrorText
	}

	switch {
	case metadata.Name == "":
		return false, apperrors.InvalidFileNameErrorText
	case metadata.IsFile && metadata.Mimetype == "":
		return false, apperrors.InvalidMimeTypeErrorText
	default:
		return true, ""
	}
}

func (u usecase) validateGetDocumentsRequest(request dto.GetDocumentsRequest) (dto.GetDocumentsRequest, bool, apperrors.ErrorText) {
	if request.Limit == 0 {
		request.Limit = 10
	}

	return request, true, ""
}
