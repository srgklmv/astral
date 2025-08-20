package user

const (
	passwordRegex             = "^[a-zA-Z0-9!&*.,#@$]{8,20}$"
	passwordUppercaseRegex    = "^.*[A-Z]{1,}.*$"
	passwordLowercaseRegex    = "^.*[a-z]{1,}.*$"
	passwordSpecialCharsRegex = "^.*[!&*.,#@$]{1,}.*$"
	passwordDigitRegex        = "^.*[0-9]{1,}.*$"
)

// Validate checks if given password meets secure requirements.
func ValidatePassword(password string) (bool, error) {
	panic("not implemented")
}

func HashPassword(password, salt string) (string, error) {
	panic("not implemented")
}

func ValidateHashedPassword(password, passwordFromDB, salt string) (bool, error) {
	panic("not implemented")
}
