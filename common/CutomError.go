package common

import "net/http"

type CustomError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func InternalServerError() *CustomError {
	return &CustomError{
		Message: "Internal Server Error",
		Code:    http.StatusInternalServerError,
	}
}

func Error(err string, errorCode int) CustomError {
	return CustomError{
		Message: err,
		Code:    errorCode,
	}
}
