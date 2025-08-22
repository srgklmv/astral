package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/srgklmv/astral/internal/models/apperrors"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
)

type documentsUsecase interface {
	UploadDocument(ctx context.Context, token string, meta dto.UploadDocumentRequestMetadata, json dto.UploadDocumentRequestJSON, file *bytes.Buffer) (dto.APIResponse[any, *dto.UploadFileResponse], int)
	// TODO: is token needed here or it can be used in middleware?
	GetDocuments(ctx context.Context, token, login, filterKey, filterValue string, limit int) (dto.APIResponse[any, *dto.GetDocumentsResponse], int)
	// TODO: What should be in headers of getting files?
	GetDocumentsHead(ctx context.Context, token, login, filterKey, filterValue string, limit int) (bool, int)
	GetDocument(ctx context.Context, token, id string) (dto.APIResponse[any, *dto.GetDocumentResponse], int)
	// TODO: What should be in headers of getting file?
	GetDocumentHead(ctx context.Context, token, id string) (bool, int)
	DeleteDocument(ctx context.Context, token, id string) (dto.APIResponse[any, *dto.DeleteDocumentResponse], int)
}

func (c controller) UploadDocument(fc *fiber.Ctx) error {
	var meta dto.UploadDocumentRequestMetadata
	metaString := fc.FormValue("meta")
	err := json.Unmarshal([]byte(metaString), &meta)
	if err != nil {
		logger.Error("request parsing error", slog.String("error", err.Error()))
		return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BodyParsingErrorCode,
			Text: apperrors.BodyParsingErrorText,
		}, nil, nil))
	}

	var jsonData dto.UploadDocumentRequestJSON
	jsonString := fc.FormValue("json")
	err = json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		logger.Error("request parsing error", slog.String("error", err.Error()))
		return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BodyParsingErrorCode,
			Text: apperrors.BodyParsingErrorText,
		}, nil, nil))
	}

	fileHeader, err := fc.FormFile("file")
	if err != nil {
		logger.Error("request parsing error", slog.String("error", err.Error()))
		return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BodyParsingErrorCode,
			Text: apperrors.BodyParsingErrorText,
		}, nil, nil))
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Error("file parsing error", slog.String("error", err.Error()))
		return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BodyParsingErrorCode,
			Text: apperrors.BodyParsingErrorText,
		}, nil, nil))
	}

	buf := bytes.NewBuffer(nil)
	n, err := io.Copy(buf, file)
	if err != nil {
		logger.Error("file parsing error", slog.String("error", err.Error()))
		return fc.Status(http.StatusInternalServerError).JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.FileUploadingErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil))
	}
	if n == 0 {
		return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BodyParsingErrorCode,
			Text: apperrors.BodyParsingErrorText,
		}, nil, nil))
	}

	token := meta.Token

	result, status := c.documentsUsecase.UploadDocument(fc.Context(), token, meta, jsonData, buf)

	return fc.Status(status).JSON(result)
}

func (c controller) GetDocuments(fc *fiber.Ctx) error {
	panic("not implemented")
}

func (c controller) GetDocumentsHead(fc *fiber.Ctx) error {
	panic("not implemented")
}

func (c controller) GetDocument(fc *fiber.Ctx) error {
	panic("not implemented")
}

func (c controller) GetDocumentHead(fc *fiber.Ctx) error {
	panic("not implemented")
}

func (c controller) DeleteDocument(fc *fiber.Ctx) error {
	panic("not implemented")
}
