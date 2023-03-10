package utils

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    uint   `json:"code"`
	Data    any    `json:"data"`
}

func BuildResponse(message string, statusCode uint, data any) Response {
	res := Response{
		Success: true,
		Message: message,
		Code:    statusCode,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, statusCode uint) Response {
	res := Response{
		Success: false,
		Message: message,
		Code:    statusCode,
		Data:    nil,
	}
	return res
}
