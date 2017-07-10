package auth

import (
	"fmt"

	"gopkg.in/ini.v1"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/registry"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/settings"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

const (
	// ComponentName is the name of this component.
	ComponentName = "app_auth"
)

var (
	Component *AppAuthComponent
)

type Configuration struct {
}

type AppAuthComponent struct {
	config   *Configuration
	registry *registry.Registry
}

// SetupIni configures lxdweb by using gopkg.in/ini.v1
func (c *AppAuthComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	cfg := &Configuration{}

	if err := iniCfg.Section(ComponentName).MapTo(cfg); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	return c.SetupStruct(cfg, configFile, debug)
}

// SetupStruct Configures lxdweb from a struct
func (c *AppAuthComponent) SetupStruct(cfg *Configuration, configFile string, debug bool) error {
	c.config = cfg

	auth.PermissionCreate(PermissionUserCreate, shared.RoleSuperAdmin, shared.RoleAdmin)
	auth.PermissionCreate(PermissionUserDelete, shared.RoleSuperAdmin, shared.RoleAdmin)
	auth.PermissionCreate(
		PermissionUpdateProfile,
		shared.RoleSuperAdmin, shared.RoleAdmin, shared.RoleUser, shared.RoleViewer,
	)
	auth.PermissionCreate(PermissionEditSettings, shared.RoleSuperAdmin, shared.RoleAdmin)

	authConfig := auth.Component.ConfigGet()
	sdb, apiErr := settings.Component.DBGet(nil)
	if apiErr != nil {
		return apiErr.Reason
	}

	s := &shared.Property{
		Key:            SettingAllowSignup,
		Value:          authConfig.AllowSignup,
		Type:           shared.PropertyTypeBool,
		PermissionView: auth.PermissionGuest,
		PermissionEdit: PermissionEditSettings,
	}
	sdb.SettingCreate(s, false)

	s = &shared.Property{
		Key:            SettingGoogleClientID,
		Value:          authConfig.GoogleClientID,
		Type:           shared.PropertyTypeString,
		PermissionView: auth.PermissionGuest,
		PermissionEdit: PermissionReadonlySetting,
	}
	sdb.SettingCreate(s, true)

	s = &shared.Property{
		Key:            SettingFacebookClientID,
		Value:          authConfig.FacebookClientID,
		Type:           shared.PropertyTypeString,
		PermissionView: auth.PermissionGuest,
		PermissionEdit: PermissionReadonlySetting,
	}
	sdb.SettingCreate(s, true)

	s = &shared.Property{
		Key:            SettingGithubClientID,
		Value:          authConfig.GithubClientID,
		Type:           shared.PropertyTypeString,
		PermissionView: auth.PermissionGuest,
		PermissionEdit: PermissionReadonlySetting,
	}
	sdb.SettingCreate(s, true)

	s = &shared.Property{
		Key:            SettingTwitterEnabled,
		Value:          authConfig.TwitterKey != "",
		Type:           shared.PropertyTypeBool,
		PermissionView: auth.PermissionGuest,
		PermissionEdit: PermissionReadonlySetting,
	}
	sdb.SettingCreate(s, true)

	return nil
}

// SetupEcho sets an echo instance up.
func (c *AppAuthComponent) SetupEcho(e *echo.Echo) error {
	apiRegisterRoutes(e)
	return nil
}

func (c *AppAuthComponent) Shutdown() error {
	return nil
}

func (c *AppAuthComponent) RegistrySet(r *registry.Registry) {
	c.registry = r
}

func (c *AppAuthComponent) NameGet() string {
	return ComponentName
}

func (c *AppAuthComponent) WeightGet() int {
	return 1001
}

func init() {
	Component = &AppAuthComponent{}
	registry.Instance().Register(ComponentName, Component)
}
