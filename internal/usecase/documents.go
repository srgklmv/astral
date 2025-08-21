package usecase

import (
	"context"

	"github.com/srgklmv/astral/internal/models/dto"
)

func (u usecase) UploadDocument(ctx context.Context, data dto.UploadDocumentRequest) dto.APIResponse[any, dto.UploadFileResponse] {
	panic("not implemented")
}

func (u usecase) GetDocuments(ctx context.Context, token, login, filterKey, filterValue string, limit int) dto.APIResponse[any, dto.GetDocumentsResponse] {
	panic("not implemented")
}

func (u usecase) GetDocumentsHead(ctx context.Context, token, login, filterKey, filterValue string, limit int) bool {
	panic("not implemented")
}

func (u usecase) GetDocument(ctx context.Context, token, id string) dto.APIResponse[any, dto.GetDocumentResponse] {
	panic("not implemented")
}

func (u usecase) GetDocumentHead(ctx context.Context, token, id string) bool {
	panic("not implemented")
}

func (u usecase) DeleteDocument(ctx context.Context, token, id string) dto.APIResponse[any, dto.DeleteDocumentResponse] {
	panic("not implemented")
}
