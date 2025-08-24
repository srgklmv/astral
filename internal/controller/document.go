package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
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
	GetDocument(ctx context.Context, token, id string) (dto.APIResponse[any, any], []byte, map[string]string, int)
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

	fileHeader, _ := fc.FormFile("file")

	var file multipart.File
	if fileHeader != nil {
		file, err = fileHeader.Open()
		if err != nil {
			logger.Error("file parsing error", slog.String("error", err.Error()))
			return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
				Code: apperrors.BodyParsingErrorCode,
				Text: apperrors.BodyParsingErrorText,
			}, nil, nil))
		}

		defer func() {
			if err := file.Close(); err != nil {
				logger.Error("file closing error", slog.String("error", err.Error()))
			}
		}()
	}

	buf := bytes.NewBuffer(nil)
	if file != nil {
		n, err := io.Copy(buf, file)
		if err != nil {
			logger.Error("file parsing error", slog.String("error", err.Error()))
			return fc.Status(http.StatusInternalServerError).JSON(dto.NewAPIResponse[any, any](&dto.Error{
				Code: apperrors.DocumentUploadingErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil))
		}
		if n == 0 {
			return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
				Code: apperrors.BodyParsingErrorCode,
				Text: apperrors.BodyParsingErrorText,
			}, nil, nil))
		}
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
	var req dto.GetDocumentRequest
	_ = fc.BodyParser(&req)

	id := fc.Params("id")

	response, file, headers, status := c.documentsUsecase.GetDocument(fc.Context(), req.Token, id)

	if len(file) == 0 {
		return fc.Status(status).JSON(response)
	}

	fmt.Print("\n REMOVE ME! ", "headers: ", headers, "\n")
	// Form data must be a choice.

	// TODO: Add headers. Which exactly? Multipart? What is going on?
	// TODO: Get headers from usecase with struct.
	fc.Attachment("popa.pdf")

	return fc.Status(status).Send(file)
}

func (c controller) GetDocumentHead(fc *fiber.Ctx) error {
	panic("not implemented")
}

func (c controller) DeleteDocument(fc *fiber.Ctx) error {
	documentID := fc.Params("id")

	var request dto.DeleteDocumentRequest
	err := fc.BodyParser(&request)
	if err != nil {
		return fc.Status(http.StatusBadRequest).JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.BodyParsingErrorText,
		}, nil, nil))
	}

	response, status := c.documentsUsecase.DeleteDocument(fc.Context(), request.Token, documentID)

	return fc.Status(status).JSON(response)
}
