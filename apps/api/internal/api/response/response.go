package response

type ErrorResponse struct {
	Code    string `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewError(code, errorStatus, msg string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Error:   errorStatus,
		Message: msg,
	}
}
