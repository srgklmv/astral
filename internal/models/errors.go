package models

type (
	ErrorText string
	ErrorCode int
)

const (
	BodyParsingErrorText ErrorText = "Bad request."
)

const (
	BodyParsingErrorCode ErrorCode = 80085
)
