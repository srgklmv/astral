package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/srgklmv/astral/internal/models"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
)

type authUsecase interface {
	Register(token, login, password string) (dto.APIResponse[dto.RegisterResponse, any], int)
	Auth(login, password string) (dto.APIResponse[dto.AuthResponse, any], int)
	Logout(token string) (dto.APIResponse[dto.LogoutResponse, any], int)
}

func (c controller) Register(fc *fiber.Ctx) error {
	var request dto.RegisterRequest
	err := fc.BodyParser(&request)
	if err != nil {
		logger.Error("request parsing error", slog.String("error", err.Error()))
		return fc.JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: models.BodyParsingErrorCode,
			Text: models.BodyParsingErrorText,
		}, nil, nil))
	}

	result, status := c.authUsecase.Register(request.Token, request.Login, request.Password)

	return fc.Status(status).JSON(result)
}

func (c controller) Auth(fc *fiber.Ctx) error {
	var request dto.AuthRequest
	err := fc.BodyParser(&request)
	if err != nil {
		logger.Error("request parsing error", slog.String("error", err.Error()))
		return fc.JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: models.BodyParsingErrorCode,
			Text: models.BodyParsingErrorText,
		}, nil, nil))
	}

	result, status := c.authUsecase.Auth(request.Login, request.Password)

	return fc.Status(status).JSON(result)
}

func (c controller) Logout(fc *fiber.Ctx) error {
	var request dto.LogoutRequest
	err := fc.ParamsParser(&request)
	if err != nil {
		logger.Error("request parsing error", slog.String("error", err.Error()))
		return fc.JSON(dto.NewAPIResponse[any, any](&dto.Error{
			Code: models.BodyParsingErrorCode,
			Text: models.BodyParsingErrorText,
		}, nil, nil))
	}

	result, status := c.authUsecase.Logout(request.Token)

	return fc.Status(status).JSON(result)
}
