package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/srgklmv/astral/internal/models"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
)

type authUsecase interface {
	Register(token, login, password string) dto.APIResponse[dto.RegisterResponse, any]
	Auth(login, password string) dto.APIResponse[dto.AuthResponse, any]
	Logout(token string) dto.APIResponse[dto.LogoutResponse, any]
}

func (c controller) Register(fc *fiber.Ctx) error {
	var request dto.RegisterRequest
	err := fc.BodyParser(&request)
	if err != nil {
		logger.Error("request parsing error", slog.String("error", err.Error()))
		return fc.JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: models.SomeInternalErrorCode,
			Text: models.SomeInternalErrorText,
		}, nil, nil))
	}

	result := c.authUsecase.Register(request.Token, request.Login, request.Password)

	return fc.JSON(result)
}

func (c controller) Auth(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}

func (c controller) Logout(fc *fiber.Ctx) error {
	panic("not implemented")
	return nil
}
