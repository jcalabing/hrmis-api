package errors

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func NewErrorResponse(code int, message string, details interface{}) *ErrorResponse {

	if message == "" {
		message = "An error occurred."
	}
	if details == "" {
		details = "Bad Request."
	}

	return &ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}
