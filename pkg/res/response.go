package res

type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   error       `json:"error,omitempty"`
}

func Success(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Data: data,
	}
}

func Fail(message string, data interface{}) *ErrorResponse {
	return &ErrorResponse{
		Data:    data,
		Message: message,
	}
}
