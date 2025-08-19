package models

type (
	ErrorText string
	ErrorCode int
)

const (
	SomeInternalErrorText ErrorText = "Oh no! So sad :("
)

const (
	SomeInternalErrorCode ErrorCode = 80085
)
