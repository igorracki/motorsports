package api

import (
	"errors"
	"net/http"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, context echo.Context) {
	if context.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "An internal server error occurred"
	errorType := "internal_error"

	// Check for domain errors
	if errors.Is(err, models.ErrNotFound) {
		code = http.StatusNotFound
		message = err.Error()
		errorType = "not_found"
	} else if errors.Is(err, models.ErrUnauthorized) {
		code = http.StatusUnauthorized
		message = err.Error()
		errorType = "unauthorized"
	} else if errors.Is(err, models.ErrForbidden) {
		code = http.StatusForbidden
		message = err.Error()
		errorType = "forbidden"
	} else if errors.Is(err, models.ErrInvalidInput) {
		code = http.StatusBadRequest
		message = err.Error()
		errorType = "invalid_input"
	} else if errors.Is(err, models.ErrConflict) {
		code = http.StatusConflict
		message = err.Error()
		errorType = "conflict"
	} else if echoError, ok := err.(*echo.HTTPError); ok {
		code = echoError.Code
		if msg, ok := echoError.Message.(string); ok {
			message = msg
		}
		errorType = "http_error"
	}

	response := models.ErrorResponse{
		Error:   errorType,
		Message: message,
	}

	if err := context.JSON(code, response); err != nil {
		context.Logger().Error(err)
	}
}
