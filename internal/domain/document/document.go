package document

import (
	"time"
)

type Document struct {
	ID        int
	Filename  string
	IsPublic  bool
	IsFile    bool
	Mimetype  string
	GrantedTo []string
	CreatedAt time.Time
	JSON      map[string]any
	File      []byte
}
