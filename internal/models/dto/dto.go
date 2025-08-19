package dto

import (
	"github.com/srgklmv/astral/internal/models"
)

// TODO: Check init values Error struct in APIResponse.
type APIResponse[Response, Data any] struct {
	Error    *Error   `json:"error,omitempty"`
	Response Response `json:"response,omitempty"`
	Data     Data     `json:"data,omitempty"`
}

type Error struct {
	Code models.ErrorCode `json:"code"`
	Text models.ErrorText `json:"text"`
}

// TODO: Check if nil Error value is ok not to return it in response.
// TODO: Check if nil on response or data is valid (no runtime errors).
func NewAPIResponse[Response, Data any](error *Error, response Response, data Data) APIResponse[Response, Data] {
	return APIResponse[Response, Data]{
		Error:    error,
		Response: response,
		Data:     data,
	}
}
