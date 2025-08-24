package dto

import (
	"time"

	"github.com/srgklmv/astral/internal/domain/document"
)

type (
	UploadDocumentRequest struct {
		Metadata string                    `form:"meta"`
		JSON     string                    `form:"json"`
		File     UploadDocumentRequestFile `form:"file"`
	}
	UploadDocumentRequestMetadata struct {
		Name     string `json:"name"`
		IsFile   bool   `json:"file"`
		IsPublic bool   `json:"public"`
		// Token is for what?
		Token     string   `json:"token"`
		Mimetype  string   `json:"mime"`
		GrantedTo []string `json:"grant"`
	}
	UploadDocumentRequestJSON map[string]any
	UploadDocumentRequestFile []byte
)

type UploadFileResponse struct {
	JSON     map[string]any `json:"json,omitempty"`
	Filename string         `json:"file"`
}

type (
	GetDocumentsRequest struct {
		Token string `json:"token"`
		Login string `json:"login"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Limit int    `json:"limit"`
	}
	GetDocumentsResponse struct {
		DocumentsData []DocumentData `json:"docs"`
	}
)

func NewGetDocumentsResponse() GetDocumentsResponse {
	return GetDocumentsResponse{
		DocumentsData: make([]DocumentData, 0),
	}
}

func (r GetDocumentsResponse) FromDomain(data document.DocumentsData) GetDocumentsResponse {
	for _, v := range data {
		r.DocumentsData = append(r.DocumentsData, DocumentData{
			ID:        v.ID.String(),
			Name:      v.Filename,
			IsFile:    v.IsFile,
			IsPublic:  v.IsPublic,
			Mimetype:  v.Mimetype,
			CreatedAt: v.CreatedAt.Format(time.DateTime),
			GrantedTo: v.GrantedTo,
		})
	}

	return r
}

type DocumentData struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	IsFile    bool     `json:"file"`
	IsPublic  bool     `json:"public"`
	Mimetype  string   `json:"mime,omitempty"`
	CreatedAt string   `json:"created"`
	GrantedTo []string `json:"grant,omitempty"`
}

type (
	GetDocumentRequest struct {
		Token string `json:"token"`
	}
	GetDocumentResponse any
)

type (
	DeleteDocumentRequest struct {
		Token string `json:"token"`
	}
	DeleteDocumentResponse map[string]bool
)
