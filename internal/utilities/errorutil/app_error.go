package errorutil

import "net/http"

type Code int

const (
	CodeInternalServerError Code = 10001
	CodeBadRequest          Code = 10002
	CodeUnauthorized        Code = 10003
	CodeNotFound            Code = 10004
	CodeValidation          Code = 10005
	CodeTelegramNotLinked   Code = 10006
)

type AppError struct {
	Code       Code
	Message    string
	HTTPStatus int
}

func New(code Code) *AppError {
	return &AppError{
		Code:       code,
		Message:    defaultMessage(code),
		HTTPStatus: defaultHTTPStatus(code),
	}
}

func NewWithMessage(code Code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: defaultHTTPStatus(code),
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func defaultMessage(code Code) string {
	switch code {
	case CodeBadRequest:
		return "bad request"
	case CodeUnauthorized:
		return "unauthorized"
	case CodeNotFound:
		return "not found"
	case CodeValidation:
		return "validation error"
	case CodeTelegramNotLinked:
		return "telegram account is not linked"
	default:
		return "internal server error"
	}
}

func defaultHTTPStatus(code Code) int {
	switch code {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeNotFound:
		return http.StatusNotFound
	case CodeValidation:
		return http.StatusUnprocessableEntity
	case CodeTelegramNotLinked:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
