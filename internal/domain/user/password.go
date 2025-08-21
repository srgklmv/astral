package user

import (
	"log/slog"
	"regexp"

	"github.com/srgklmv/astral/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordRegex             = "^[a-zA-Z0-9!&*.,#@$]{8,20}$"
	passwordUppercaseRegex    = "^.*[A-Z]{1,}.*$"
	passwordLowercaseRegex    = "^.*[a-z]{1,}.*$"
	passwordSpecialCharsRegex = "^.*[!&*.,#@$]{1,}.*$"
	passwordDigitRegex        = "^.*[0-9]{1,}.*$"
)

// Validate checks if given password meets secure requirements.
func ValidatePassword(password string) (bool, error) {
	matched1, err := regexp.MatchString(passwordRegex, password)
	if err != nil {
		logger.Error("regex compilation error", slog.String("error", err.Error()))
		return false, err
	}

	matched2, err := regexp.MatchString(passwordUppercaseRegex, password)
	if err != nil {
		logger.Error("regex compilation error", slog.String("error", err.Error()))
		return false, err
	}

	matched3, err := regexp.MatchString(passwordLowercaseRegex, password)
	if err != nil {
		logger.Error("regex compilation error", slog.String("error", err.Error()))
		return false, err
	}

	matched4, err := regexp.MatchString(passwordSpecialCharsRegex, password)
	if err != nil {
		logger.Error("regex compilation error", slog.String("error", err.Error()))
		return false, err
	}

	matched5, err := regexp.MatchString(passwordDigitRegex, password)
	if err != nil {
		logger.Error("regex compilation error", slog.String("error", err.Error()))
		return false, err
	}

	return matched1 && matched2 && matched3 && matched4 && matched5, nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("password hashing error", slog.String("error", err.Error()))
		return "", err
	}

	return string(hashed), nil
}

func IsValidPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
