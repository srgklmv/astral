package usecase

import (
	"github.com/srgklmv/astral/internal/domain/user"
)

type repository interface {
	userRepository
}

type userRepository interface {
	IsLoginExists(login string) (bool, error)
	IsAdminTokenValid(token string) (bool, error)
	CreateUser(login, hashedPassword string, isAdmin bool) (user.User, error)
	GetByLogin(login string) (user.User, error)
	ValidatePassword(userID int, hashedPassword string) (bool, error)
	SaveAuthToken(userID int, token string) error
	DeleteToken(login string) (bool, error)
}

type usecase struct {
	userRepository userRepository
}

func New(repository repository) *usecase {
	return &usecase{
		userRepository: repository,
	}
}
