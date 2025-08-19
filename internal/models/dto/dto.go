package dto

// TODO: Check init values Error struct in APIResponse.
type APIResponse[Response, Data any] struct {
	Error    Error    `json:"error,omitempty"`
	Response Response `json:"response,omitempty"`
	Data     Data     `json:"data,omitempty"`
}

type Error struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// TODO: Check if nil Error value is ok not to return it in response.
// TODO: Check if nil on response or data is valid (no runtime errors).
func NewAPIResponse[Response, Data any](error Error, response Response, data Data) APIResponse[Response, Data] {
	return APIResponse[Response, Data]{
		Error:    error,
		Response: response,
		Data:     data,
	}
}
