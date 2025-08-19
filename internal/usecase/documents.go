package usecase

import (
	"github.com/srgklmv/astral/internal/models/dto"
)

func (u usecase) UploadDocument(dto.UploadDocumentRequest) dto.APIResponse[any, dto.UploadFileResponse] {
	panic("not implemented")
}

func (u usecase) GetDocuments(token, login, filterKey, filterValue string, limit int) dto.APIResponse[any, dto.GetDocumentsResponse] {
	panic("not implemented")
}

func (u usecase) GetDocumentsHead(token, login, filterKey, filterValue string, limit int) bool {
	panic("not implemented")
}

func (u usecase) GetDocument(token, id string) dto.APIResponse[any, dto.GetDocumentResponse] {
	panic("not implemented")
}

func (u usecase) GetDocumentHead(token, id string) bool {
	panic("not implemented")
}

func (u usecase) DeleteDocument(token, id string) dto.APIResponse[any, dto.DeleteDocumentResponse] {
	panic("not implemented")
}
