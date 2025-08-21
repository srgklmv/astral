package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/srgklmv/astral/internal/models/dto"
)

type documentsUsecase interface {
	UploadDocument(ctx context.Context, data dto.UploadDocumentRequest) dto.APIResponse[any, dto.UploadFileResponse]
	// TODO: is token needed here or it can be used in middleware?
	GetDocuments(ctx context.Context, token, login, filterKey, filterValue string, limit int) dto.APIResponse[any, dto.GetDocumentsResponse]
	// TODO: What should be in headers of getting files?
	GetDocumentsHead(ctx context.Context, token, login, filterKey, filterValue string, limit int) bool
	GetDocument(ctx context.Context, token, id string) dto.APIResponse[any, dto.GetDocumentResponse]
	// TODO: What should be in headers of getting file?
	GetDocumentHead(ctx context.Context, token, id string) bool
	DeleteDocument(ctx context.Context, token, id string) dto.APIResponse[any, dto.DeleteDocumentResponse]
}

func (c controller) UploadDocument(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) GetDocuments(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) GetDocumentsHead(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) GetDocument(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) GetDocumentHead(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) DeleteDocument(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}
