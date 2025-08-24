package dto

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

type GetDocumentsResponse struct {
	Documents []Document `json:"docs"`
}

// TODO: Add bytes here.
type Document struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	IsFile    bool     `json:"file"`
	IsPublic  bool     `json:"public"`
	Mimetype  string   `json:"mime"`
	CreatedAt string   `json:"created"`
	GrantedTo []string `json:"grant"`
}

// TODO: Is it possible, to avoid any here? JSON of file returned.
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
