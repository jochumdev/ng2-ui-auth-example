package shared

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// APIError is error type that knows if its Internal
// and the HTTP StatusCode
type APIError struct {
	Reason   error
	Internal bool
	Code     int
}

// APIHandleError handles Errors from API Functions,
// Preferable an APIError but takes error also.
func APIHandleError(c echo.Context, err interface{}) error {
	var errAPI *APIError

	switch err.(type) {
	case error:
		errorErr := err.(error)
		log.Printf("API ERROR: %s", errorErr.Error())
		c.JSON(http.StatusInternalServerError,
			map[string]string{"message": errorErr.Error()},
		)
		return errorErr

	case *APIError:
		errAPI = err.(*APIError)
	case APIError:
		myErr := err.(APIError)
		errAPI = &myErr
	default:
		log.Printf("ERROR: Unknown error %v reveived.", err)
		c.JSON(http.StatusInternalServerError,
			map[string]string{"message": "Unknown error received!"},
		)

		return fmt.Errorf("Unknown error %v received.", err)
	}

	if errAPI != nil {
		if errAPI.Internal {
			log.Printf("API ERROR: %s", errAPI.Reason.Error())

			if errAPI.Code != 0 {
				c.JSON(errAPI.Code,
					map[string]string{"message": "Internal server error"},
				)
			} else {
				c.JSON(http.StatusInternalServerError,
					map[string]string{"message": "Internal server error"},
				)
			}

			return errAPI.Reason
		}

		if errAPI.Code != 0 {
			log.Printf("API ERROR: %s", errAPI.Reason.Error())
			c.JSON(errAPI.Code, map[string]string{"message": errAPI.Reason.Error()})

			return errAPI.Reason
		}

		log.Printf("API ERROR: %s", errAPI.Reason.Error())
		c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": errAPI.Reason.Error()},
		)

		return errAPI.Reason
	}

	return ErrInternalServer
}
