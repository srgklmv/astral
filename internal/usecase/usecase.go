package usecase

import (
	"bytes"
	"context"

	"github.com/google/uuid"
	"github.com/srgklmv/astral/internal/domain/document"
	"github.com/srgklmv/astral/internal/domain/user"
)

type repository interface {
	userRepository
	documentRepository
}

type documentRepository interface {
	UploadDocument(ctx context.Context, login, filename string, isFile bool, mimetype string, isPublic bool, grantedTo []string, json map[string]any, file *bytes.Buffer) (document.Document, error)
	DeleteDocument(ctx context.Context, id uuid.UUID) error
	GetDocument(ctx context.Context, id uuid.UUID) (document.Document, error)
}

type userRepository interface {
	IsLoginExists(ctx context.Context, login string) (bool, error)
	IsAdminTokenValid(ctx context.Context, token string) (bool, error)
	CreateUser(ctx context.Context, login, hashedPassword string, isAdmin bool) (user.User, error)
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
	ValidatePassword(ctx context.Context, userID int, hashedPassword string) (bool, error)
	SaveAuthToken(ctx context.Context, login, token string) error
	DeleteToken(ctx context.Context, token string) error
	GetUserHashedPassword(ctx context.Context, login string) (string, error)
	DeleteAllUserTokens(ctx context.Context, login string) error
	IsAuthTokenExists(ctx context.Context, token string) (bool, error)
	GetUserLoginByAuthToken(ctx context.Context, token string) (string, error)
}

type usecase struct {
	userRepository     userRepository
	documentRepository documentRepository
}

func New(repository repository) *usecase {
	return &usecase{
		userRepository:     repository,
		documentRepository: repository,
	}
}
