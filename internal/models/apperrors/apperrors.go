package apperrors

type (
	ErrorText string
	ErrorCode int
)

// Global blocks.
const (
	InternalErrorText    ErrorText = "Internal error. Chat support, please."
	BodyParsingErrorText ErrorText = "Bad request."
)

const (
	InternalErrorErrorCode ErrorCode = 80085
	BodyParsingErrorCode   ErrorCode = 100 + iota
	RepositoryCallErrorCode
	RegexErrorCode
	BadRequestErrorCode
)

// Auth blocks.
const (
	RegisterBadLoginErrorText     ErrorText = "Login length must be between 8 and 20."
	RegisterLoginTakenErrorText   ErrorText = "Login already taken."
	RegisterBadPasswordErrorText  ErrorText = "Password length must be between 8 and 20, contains at least one upper and one lower case letter, one digit and one special symbol (!&*.,#@$)."
	AuthWrongCredentialsErrorText ErrorText = "Wrong credentials."
	UnauthorizedErrorText         ErrorText = "Unauthorized."
)

const (
	PasswordHashErrorCode ErrorCode = 200 + iota
	AuthWrongLoginErrorCode
	AuthTokenGenerationErrorCode
	UnauthorizedErrorCode
)

// Document blocks.
const (
	InvalidFileNameErrorText ErrorText = "Invalid file name."
	InvalidMimeTypeErrorText ErrorText = "Invalid mime type."
	EmptyFileErrorText       ErrorText = "Empty file."
)

const (
	FileUploadingErrorCode ErrorCode = 300 + iota
)
