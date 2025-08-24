package document

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	Data
	JSON map[string]any
	File []byte
}

type Data struct {
	ID        uuid.UUID
	Filename  string
	IsPublic  bool
	IsFile    bool
	Mimetype  string
	GrantedTo []string
	CreatedAt time.Time
	Owner     string
}

type DocumentsData []Data
