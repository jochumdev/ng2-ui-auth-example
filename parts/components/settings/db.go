package settings

import (
	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/shared"
)

type DBSettingsApi interface {
	ConfigSet(c *Configuration)
	DBGet() (DBSettingsApi, *shared.APIError)
	DBFromContext(c echo.Context) (DBSettingsApi, *shared.APIError)

	/**
	 * Never call the methods below without DBGet() or DBFromContext()
	 * as some drivers need to clone itself or get the real DBDriver.
	 **/
	Initialize() error

	// You shouldn't call SettingsMap to often, for most drivers this is
	// a very expensive call.
	SettingsList() ([]*shared.Property, *shared.APIError)
	SettingsMap() (map[string]*shared.Property, *shared.APIError)
	SettingGet(string) (*shared.Property, *shared.APIError)
	SettingHas(*shared.Property) (bool, *shared.APIError)
	SettingCreate(*shared.Property, bool) *shared.APIError
	SettingUpdate(*shared.Property) *shared.APIError
}
