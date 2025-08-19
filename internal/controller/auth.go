package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/srgklmv/astral/internal/models/dto"
)

type authUsecase interface {
	Register(token, login, password string) dto.APIResponse[dto.RegisterResponse, any]
	Auth(login, password string) dto.APIResponse[dto.AuthResponse, any]
	Logout(token string) dto.APIResponse[dto.LogoutResponse, any]
}

func (c controller) Register(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) Auth(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) Logout(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}
