package usecase

import (
	"context"
	"log/slog"
	"net/http"

	userDomain "github.com/srgklmv/astral/internal/domain/user"
	"github.com/srgklmv/astral/internal/models/apperrors"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
)

func (u usecase) Register(ctx context.Context, token, login, password string) (dto.APIResponse[*dto.RegisterResponse, any], int) {
	matched, err := userDomain.ValidateLogin(login)
	if err != nil {
		logger.Error("login validation error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.RegexErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}
	if !matched {
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.BadRequestErrorCode,
				Text: apperrors.RegisterBadLoginErrorText,
			}, nil, nil,
		), http.StatusBadRequest
	}

	matched, err = userDomain.ValidatePassword(password)
	if err != nil {
		logger.Error("password validation error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.RegexErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}
	if !matched {
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.BadRequestErrorCode,
				Text: apperrors.RegisterBadPasswordErrorText,
			}, nil, nil,
		), http.StatusBadRequest
	}

	isLoginTaken, err := u.userRepository.IsLoginExists(ctx, login)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.RepositoryCallErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}
	if isLoginTaken {
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.BadRequestErrorCode,
				Text: apperrors.RegisterLoginTakenErrorText,
			}, nil, nil,
		), http.StatusBadRequest
	}

	var isAdmin bool
	if token != "" {
		isAdmin, err = u.userRepository.IsAdminTokenValid(ctx, token)
		if err != nil {
			logger.Error("repository call error", slog.String("error", err.Error()))
			return dto.NewAPIResponse[*dto.RegisterResponse, any](
				&dto.Error{
					Code: apperrors.RepositoryCallErrorCode,
					Text: apperrors.InternalErrorText,
				}, nil, nil,
			), http.StatusInternalServerError
		}
	}

	hashedPassword, err := userDomain.HashPassword(password)
	if err != nil {
		logger.Error("password hashing error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.PasswordHashErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}

	user, err := u.userRepository.CreateUser(ctx, login, hashedPassword, isAdmin)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: apperrors.RepositoryCallErrorCode,
				Text: apperrors.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}

	return dto.NewAPIResponse[*dto.RegisterResponse, any](
		nil,
		&dto.RegisterResponse{Login: user.Login},
		nil,
	), http.StatusCreated
}

func (u usecase) Auth(ctx context.Context, login, password string) (dto.APIResponse[*dto.AuthResponse, any], int) {
	if login == "" || password == "" {
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.AuthWrongCredentialsErrorText,
		}, nil, nil), http.StatusBadRequest
	}

	userExists, err := u.userRepository.IsLoginExists(ctx, login)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}
	if !userExists {
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.AuthWrongCredentialsErrorText,
		}, nil, nil), http.StatusBadRequest
	}

	hashed, err := u.userRepository.GetUserHashedPassword(ctx, login)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}
	if hashed == "" {
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.AuthWrongCredentialsErrorText,
		}, nil, nil), http.StatusBadRequest
	}

	if !userDomain.IsValidPassword(password, hashed) {
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.BadRequestErrorCode,
			Text: apperrors.AuthWrongCredentialsErrorText,
		}, nil, nil), http.StatusBadRequest
	}

	err = u.userRepository.DeleteAllUserTokens(ctx, login)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	// TODO: Add caching.

	token := userDomain.GenerateAuthToken()

	err = u.userRepository.SaveAuthToken(ctx, login, token)

	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.AuthResponse, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	return dto.NewAPIResponse[*dto.AuthResponse, any](
		nil,
		&dto.AuthResponse{Token: token},
		nil,
	), http.StatusCreated
}

func (u usecase) Logout(ctx context.Context, token string) (dto.APIResponse[*dto.LogoutResponse, any], int) {
	// TODO: Add comparison to header token.
	err := u.userRepository.DeleteToken(ctx, token)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.LogoutResponse, any](&dto.Error{
			Code: apperrors.RepositoryCallErrorCode,
			Text: apperrors.InternalErrorText,
		}, nil, nil), http.StatusInternalServerError
	}

	return dto.NewAPIResponse[*dto.LogoutResponse, any](nil, &dto.LogoutResponse{
		token: true,
	}, nil), http.StatusOK
}

func (u usecase) validateAuthToken(ctx context.Context, token string) (isValid bool, login string, err error) {
	if token == "" {
		return false, "", nil
	}

	login, err = u.userRepository.GetUserLoginByAuthToken(ctx, token)
	if err != nil {
		logger.Error("repository call error", slog.String("error", err.Error()))
		return false, "", err
	}

	return login != "", login, nil
}
