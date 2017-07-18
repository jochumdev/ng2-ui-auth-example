package auth

import (
	"github.com/labstack/echo"

	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

type DBAuthAPI interface {
	ConfigSet(c *Configuration)
	SetupEcho(e *echo.Echo) error

	DBGet() (DBAuthAPI, *shared.APIError)
	DBFromContext(c echo.Context) (DBAuthAPI, *shared.APIError)

	/**
	 * Never call the methods below without DBGet() or DBFromContext()
	 * as some drivers need to clone itself or get the real DBDriver.
	 **/
	Initialize() error

	UserFindByID(id string) (*User, *shared.APIError)
	UserFindByUsername(username string) (*User, *shared.APIError)
	UserFindByProperty(property string, value string) (*User, *shared.APIError)
	UserCreate(*User) *shared.APIError
	UserUpdate(*User) *shared.APIError
	UserDelete(username string) *shared.APIError
	UserCount() (int, *shared.APIError)

	// Permissions() ([]string, *shared.APIError)
	PermissionCreate(permission string) *shared.APIError
	//
	// Roles() ([]*Role, *shared.APIError)
	RoleFind(name string) (*Role, *shared.APIError)
	RoleCreate(*Role) *shared.APIError
	RoleUpdate(*Role) *shared.APIError
	// RoleDelete(*Role) *shared.APIError
}

func DBGet() (DBAuthAPI, *shared.APIError) {
	return Component.DriverGet().DBGet()
}

func DBFromContext(c echo.Context) (DBAuthAPI, *shared.APIError) {
	return Component.DriverGet().DBFromContext(c)
}
