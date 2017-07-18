package settings

import (
	"fmt"

	"gopkg.in/ini.v1"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/components/auth"
	"github.com/pcdummy/ng2-ui-auth-example/server/components/registry"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

const (
	// ComponentName is the name of this component.
	ComponentName = "app_settings"
)

var (
	Component *AppSettingsComponent
)

type Configuration struct {
}

type AppSettingsComponent struct {
	config   *Configuration
	registry *registry.Registry
}

// SetupIni configures lxdweb by using gopkg.in/ini.v1
func (c *AppSettingsComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	cfg := &Configuration{}

	if err := iniCfg.Section(ComponentName).MapTo(cfg); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	return c.SetupStruct(cfg, configFile, debug)
}

// SetupStruct Configures lxdweb from a struct
func (c *AppSettingsComponent) SetupStruct(cfg *Configuration, configFile string, debug bool) error {
	c.config = cfg

	auth.PermissionCreate(PermissionUpdateSetting, shared.RoleSuperAdmin, shared.RoleAdmin)

	return nil
}

// SetupEcho sets an echo instance up.
func (c *AppSettingsComponent) SetupEcho(e *echo.Echo) error {
	apiRegisterRoutes(e)
	return nil
}

func (c *AppSettingsComponent) Shutdown() error {
	return nil
}

func (c *AppSettingsComponent) RegistrySet(r *registry.Registry) {
	c.registry = r
}

func (c *AppSettingsComponent) NameGet() string {
	return ComponentName
}

func (c *AppSettingsComponent) WeightGet() int {
	return 1001
}

func init() {
	Component = &AppSettingsComponent{}
	registry.Instance().Register(ComponentName, Component)
}
