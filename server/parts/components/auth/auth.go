package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

var (
	// ErrWrongUsernameOrPassword an error that indicates the user logged in
	// with a wrong password or username.
	ErrWrongUsernameOrPassword = errors.New("Wrong username or password")

	// ErrTokenExpired the given token has expired.
	ErrTokenExpired = errors.New("Token Expired")

	// ErrUnknownDatabase the configured database is unknown.
	ErrUnknownDatabase = errors.New("Unknown Database")

	ErrPermissionExists = errors.New("Permission does exists")

	ErrRoleExists = errors.New("Role does exists")

	ErrNotFound = errors.New("Not found")

	ErrUserExists = errors.New("The given user already exists.")
)

// TokenData is the type for data stored in a token.
type TokenData struct {
	ID  string
	Iat float64
	Exp float64
}

// AuthorisationWrapper gives your echo function an user and checks permissions.
func AuthorisationWrapper(callable func(c echo.Context, u *User) error, permission string) func(echo.Context) error {
	return func(c echo.Context) error {
		user := c.Get("user").(*User)

		if permission != "" {
			if !user.PermissionHas(permission) {
				return shared.APIHandleError(c, shared.ErrUnauthorized)
			}

		}

		return callable(c, user)
	}
}

// TokenValidate validates a token and returns TokenData
func TokenValidate(token string) (*TokenData, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Always check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the key for validation
		return pubKey, nil
	})

	if err != nil || !t.Valid {
		switch err.(type) {
		case *jwt.ValidationError:
			{
				vErr := err.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					return nil, ErrTokenExpired
				default:
					return nil, vErr
				}
			}
		}
	}

	tokenData := t.Claims.(jwt.MapClaims)
	return &TokenData{tokenData["ID"].(string), tokenData["iat"].(float64), tokenData["exp"].(float64)}, nil
}

// PermissionCreate tries to create the given permission in the db,
// if it succeeds it will append the permission to all roles given.
func PermissionCreate(permission string, roles ...string) *shared.APIError {
	db, apiErr := DBGet()
	if apiErr != nil {
		return apiErr
	}

	apiErr = db.PermissionCreate(permission)
	if apiErr != nil {
		return apiErr
	}

	for _, rn := range roles {
		r, apiErr := db.RoleFind(rn)

		if apiErr != nil {
			// Silently ignore unknown roles
			continue
		}

		r.PermissionAdd(permission)
		db.RoleUpdate(r)
	}

	return nil
}
