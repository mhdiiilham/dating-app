package restful

import "net/http"

// Response struct is the general response for success and fails http call.
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

// InternalServerErrorResponse function ...
func InternalServerErrorResponse(err error) Response {
	return Response{
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
		Data:    nil,
		Error:   err,
	}
}

// BadRequestErrorResponse function ...
func BadRequestErrorResponse(err error) Response {
	return Response{
		Status:  http.StatusBadRequest,
		Message: err.Error(),
		Data:    nil,
		Error:   err,
	}
}

// SuccessResponse function ...
func SuccessResponse(data any, message string, statusCodes ...int) Response {
	statusCode := http.StatusOK
	if len(statusCodes) > 0 {
		statusCode = statusCodes[0]
	}

	return Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
		Error:   nil,
	}
}
