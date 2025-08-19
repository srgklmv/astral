package models

type (
	ErrorText string
	ErrorCode int
)

// Global texts.
const (
	InternalErrorText    ErrorText = "Internal error. Chat support, please."
	BodyParsingErrorText ErrorText = "Bad request."
)

// Auth texts.
const (
	RegisterBadLoginErrorText ErrorText = "Login length must be between 8 and 20."
)

// Global codes.
const (
	InternalErrorErrorCode ErrorCode = 80085
	BodyParsingErrorCode   ErrorCode = 100 + iota
)

const (
	RegisterBadRequestErrorCode ErrorCode = 200 + iota
)
