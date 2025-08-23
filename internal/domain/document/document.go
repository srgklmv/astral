package document

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID
	Filename  string
	IsPublic  bool
	IsFile    bool
	Mimetype  string
	GrantedTo []string
	CreatedAt time.Time
	JSON      map[string]any
	File      []byte
	Owner     string
}
