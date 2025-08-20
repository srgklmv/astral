package usecase

import (
	"log/slog"
	"net/http"

	"github.com/srgklmv/astral/internal/config"
	userDomain "github.com/srgklmv/astral/internal/domain/user"
	"github.com/srgklmv/astral/internal/models"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
)

func (u usecase) Register(token, login, password string) (dto.APIResponse[*dto.RegisterResponse, any], int) {
	matched, err := userDomain.ValidateLogin(login)
	if err != nil {
		logger.Error("login validation error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.RegexErrorCode,
				Text: models.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}
	if !matched {
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.BadRequestErrorCode,
				Text: models.RegisterBadLoginErrorText,
			}, nil, nil,
		), http.StatusBadRequest
	}

	matched, err = userDomain.ValidatePassword(password)
	if err != nil {
		logger.Error("password validation error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.RegexErrorCode,
				Text: models.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}
	if !matched {
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.BadRequestErrorCode,
				Text: models.RegisterBadPasswordErrorText,
			}, nil, nil,
		), http.StatusBadRequest
	}

	taken, err := u.userRepository.IsLoginExists(login)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.RepositoryCallErrorCode,
				Text: models.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}
	if taken {
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.BadRequestErrorCode,
				Text: models.RegisterLoginTakenErrorText,
			}, nil, nil,
		), http.StatusBadRequest
	}

	var isAdmin bool
	if token != "" {
		isAdmin, err = u.userRepository.IsAdminTokenValid(token)
		if err != nil {
			logger.Error("repository call error", slog.String("error", err.Error()))
			return dto.NewAPIResponse[*dto.RegisterResponse, any](
				&dto.Error{
					Code: models.RepositoryCallErrorCode,
					Text: models.InternalErrorText,
				}, nil, nil,
			), http.StatusInternalServerError
		}
	}

	hashedPassword, err := userDomain.HashPassword(password, config.Cfg.Modules.Auth.Salt)
	if err != nil {
		logger.Error("password hashing error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.PasswordHashErrorCode,
				Text: models.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}

	user, err := u.userRepository.CreateUser(login, hashedPassword, isAdmin)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.RepositoryCallErrorCode,
				Text: models.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}

	return dto.NewAPIResponse[*dto.RegisterResponse, any](
		nil,
		&dto.RegisterResponse{Login: user.Login},
		nil,
	), http.StatusCreated
}

func (u usecase) Auth(login, password string) (dto.APIResponse[*dto.AuthResponse, any], int) {
	user, err := u.userRepository.GetByLogin(login)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: models.RepositoryCallErrorCode,
			Text: models.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}
	if user.ID == 0 {
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: models.BadRequestErrorCode,
			Text: models.AuthWrongCredentials,
		}, nil, nil), http.StatusBadRequest
	}

	hashed, err := userDomain.HashPassword(password, config.Cfg.Modules.Auth.Salt)
	if err != nil {
		logger.Error("password hashing error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: models.PasswordHashErrorCode,
			Text: models.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	valid, err := u.userRepository.ValidatePassword(user.ID, hashed)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: models.RepositoryCallErrorCode,
			Text: models.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}
	if !valid {
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: models.BadRequestErrorCode,
			Text: models.AuthWrongCredentials,
		}, nil, nil), http.StatusBadRequest
	}

	token, err := userDomain.GenerateAuthToken(user.ID, config.Cfg.Modules.Auth.TokenSalt)
	if err != nil {
		logger.Error("token generation error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: models.AuthTokenGenerationErrorCode,
			Text: models.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	err = u.userRepository.SaveAuthToken(user.ID, token)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: models.RepositoryCallErrorCode,
			Text: models.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	return dto.NewAPIResponse[*dto.AuthResponse, any](
		nil,
		&dto.AuthResponse{Token: token},
		nil,
	), http.StatusCreated
}

func (u usecase) Logout(token string) (dto.APIResponse[dto.LogoutResponse, any], int) {

	panic("not implemented")
}
