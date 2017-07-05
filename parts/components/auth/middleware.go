package auth

import (
	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/shared"
)

const (
	bearer = "Bearer"
)

// MiddlewareJWTAuth is JSON Web Token middleware
var MiddlewareJWTAuth echo.MiddlewareFunc

// MiddlewareJWTAuthJSON is a JSON Web Token middleware that returns JSON
var MiddlewareJWTAuthJSON echo.MiddlewareFunc

func middlewareJWTAuth(returnJSON bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(bearer)

			if len(auth) < l+1 || auth[:l] != bearer {
				c.Set("user", Component.UserGuest)
				return next(c)
			}

			tdata, err := TokenValidate(auth[l+1:])
			if err != nil {
				return shared.APIHandleError(c, shared.ErrUnauthorized)
			}

			db, apiErr := DBFromContext(c)
			if err != nil {
				return shared.APIHandleError(c, *apiErr)
			}

			user, apiErr := db.UserFindByID(tdata.ID)
			if apiErr != nil {
				return shared.APIHandleError(c, *apiErr)
			}

			// Store the user in echo.Context
			c.Set("user", user)

			// Store token claims in echo.Context
			c.Set("claims", tdata)
			return next(c)
		}
	}
}
