package user

import (
	"log/slog"
	"regexp"

	"github.com/srgklmv/astral/pkg/logger"
)

const (
	loginRegex = "^[a-zA-Z0-9]{8,20}$"
)

func ValidateLogin(login string) (bool, error) {
	matched, err := regexp.MatchString(loginRegex, login)
	if err != nil {
		logger.Error("regex compilation error", slog.String("error", err.Error()))
		return false, err
	}

	return matched, nil
}
