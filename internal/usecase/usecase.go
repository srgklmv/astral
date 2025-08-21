package usecase

import (
	"context"

	"github.com/srgklmv/astral/internal/domain/user"
)

type repository interface {
	userRepository
}

type userRepository interface {
	IsLoginExists(ctx context.Context, login string) (bool, error)
	IsAdminTokenValid(ctx context.Context, token string) (bool, error)
	CreateUser(ctx context.Context, login, hashedPassword string, isAdmin bool) (user.User, error)
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
	ValidatePassword(ctx context.Context, userID int, hashedPassword string) (bool, error)
	SaveAuthToken(ctx context.Context, login, token string) error
	DeleteToken(ctx context.Context, token string) (bool, error)
	GetUserHashedPassword(ctx context.Context, login string) (string, error)
}

type usecase struct {
	userRepository userRepository
}

func New(repository repository) *usecase {
	return &usecase{
		userRepository: repository,
	}
}
