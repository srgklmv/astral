package controller

import (
	"github.com/srgklmv/astral/internal/models/dto"
)

type authUsecase interface {
	Register(token, login, password string) dto.APIResponse[dto.RegisterResponse, any]
	Auth(login, password string) dto.APIResponse[dto.AuthResponse, any]
	Logout(token string) dto.APIResponse[dto.LogoutResponse, any]
}
