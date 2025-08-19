package controller

import (
	"github.com/srgklmv/astral/internal/models/dto"
)

type docsUsecase interface {
	UploadDocument(dto.UploadDocumentRequest) dto.APIResponse[any, dto.UploadFileResponse]
	// TODO: is token needed here or it can be used in middleware?
	GetDocuments(token, login, filterKey, filterValue string, limit int) dto.APIResponse[any, dto.GetDocumentsResponse]
	// TODO: What should be in headers of getting files?
	GetDocumentsHead(token, login, filterKey, filterValue string, limit int) bool

	GetDocument(token, id string) dto.APIResponse[any, dto.GetDocumentResponse]
	// TODO: What should be in headers of getting file?
	GetDocumentHead(token, id string) bool

	DeleteDocument(token, id string) dto.APIResponse[any, dto.DeleteDocumentResponse]
}
