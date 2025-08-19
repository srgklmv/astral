package usecase

import (
	"log/slog"
	"net/http"
	"regexp"

	"github.com/srgklmv/astral/internal/models"
	"github.com/srgklmv/astral/internal/models/dto"
	"github.com/srgklmv/astral/pkg/logger"
)

const (
	loginRegex                = "^[a-zA-Z0-9]{8,20}$"
	passwordRegex             = "^[a-zA-Z0-9!&*.,#@$]{8,20}$"
	passwordUppercaseRegex    = "^.*[A-Z]{1,}.*$"
	passwordLowercaseRegex    = "^.*[a-z]{1,}.*$"
	passwordSpecialCharsRegex = "^.*[!&*.,#@$]{1,}.*$"
	passwordDigitRegex        = "^.*[0-9]{1,}.*$"
)

func (u usecase) Register(token, login, password string) (dto.APIResponse[*dto.RegisterResponse, any], int) {
	matched, err := regexp.MatchString(loginRegex, login)
	if err != nil {
		logger.Error("regex compilation error", slog.String("error", err.Error()))
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.InternalErrorErrorCode,
				Text: models.InternalErrorText,
			}, nil, nil,
		), http.StatusInternalServerError
	}
	if !matched {
		return dto.NewAPIResponse[*dto.RegisterResponse, any](
			&dto.Error{
				Code: models.InternalErrorErrorCode,
				Text: models.InternalErrorText,
			}, nil, nil,
		), http.StatusBadRequest
	}

	b := []byte(login)
	if len(b) < 8 || len(b) > 20 {
		return dto.APIResponse[dto.RegisterResponse, any]{
			Error: &dto.Error{
				Code: models.RegisterBadRequestErrorCode,
				Text: models.RegisterBadLoginErrorText,
			},
		}, http.StatusBadRequest
	}
	// TODO: Check for password to pass requirements.
	b = []byte(password)
	if len(b) < 8 || len(b) > 20 {
		return dto.APIResponse[dto.RegisterResponse, any]{
			Error: &dto.Error{
				Code: models.RegisterBadRequestErrorCode,
				Text: models.RegisterBadLoginErrorText,
			},
		}, http.StatusBadRequest
	}

	// TODO: Check if login not taken.
	// TODO: Check for admin token.

	panic("not implemented")
}

func (u usecase) Auth(login, password string) (dto.APIResponse[dto.AuthResponse, any], int) {
	panic("not implemented")
}

func (u usecase) Logout(token string) (dto.APIResponse[dto.LogoutResponse, any], int) {
	panic("not implemented")
}
