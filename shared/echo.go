package shared

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

var (
	ErrInternalServer = echo.NewHTTPError(http.StatusInternalServerError)
	ErrBadRequest     = echo.NewHTTPError(http.StatusBadRequest)
	ErrUnauthorized   = &APIError{
		Reason: errors.New("Unauthorized"),
		Code:   http.StatusUnauthorized,
	}
)
