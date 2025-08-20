package repository

import (
	"github.com/srgklmv/astral/internal/domain/user"
)

func (r repository) IsLoginExists(login string) (bool, error) {
	panic("not implemented")
}
func (r repository) IsAdminTokenValid(token string) (bool, error) {
	panic("not implemented")
}
func (r repository) CreateUser(login, hashedPassword string, isAdmin bool) (user.User, error) {
	panic("not implemented")
}
func (r repository) GetByLogin(login string) (user.User, error) {
	panic("not implemented")
}
func (r repository) ValidatePassword(userID int, hashedPassword string) (bool, error) {
	panic("not implemented")
}
func (r repository) SaveAuthToken(userID int, token string) error {
	panic("not implemented")
}
func (r repository) DeleteToken(login string) (bool, error) {
	panic("not implemented")
}
