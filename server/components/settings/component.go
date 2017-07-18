package settings

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/components/registry"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

const (
	// ComponentName is the name of this component.
	ComponentName = "com_settings"
)

var (
	Component *SettingsComponent
)

type Configuration struct {
	Debug          bool
	SettingsDBType string
	SettingsDBUrl  string
}

type SettingsComponent struct {
	config   *Configuration
	registry *registry.Registry

	driver  DBSettingsApi
	drivers map[string](DBSettingsApi)
}

// SetupIni configures lxdweb by using gopkg.in/ini.v1
func (c *SettingsComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	cfg := &Configuration{}

	if err := iniCfg.Section(ComponentName).MapTo(cfg); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	return c.SetupStruct(cfg, configFile, debug)
}

// SetupStruct Configures lxdweb from a struct
func (c *SettingsComponent) SetupStruct(cfg *Configuration, configFile string, debug bool) error {
	c.config = cfg
	c.config.Debug = debug

	var (
		driver DBSettingsApi
		ok     bool
	)

	if driver, ok = c.drivers[c.config.SettingsDBType]; !ok {
		log.Fatalf(
			"Unknown db backend '%s' configured for settings.",
			c.config.SettingsDBType,
		)
	}

	c.driver = driver
	c.driver.ConfigSet(c.config)
	db, apiErr := c.driver.DBGet()
	if apiErr != nil {
		log.Fatalf("Failed to get the db: %v", apiErr.Reason)
	}
	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize the settings db: %v", err)
	}

	return nil
}

// SetupEcho sets an echo instance up.
func (c *SettingsComponent) SetupEcho(e *echo.Echo) error {
	return nil
}

func (c *SettingsComponent) Shutdown() error {
	return nil
}

func (c *SettingsComponent) RegistrySet(r *registry.Registry) {
	c.registry = r
}

func (c *SettingsComponent) NameGet() string {
	return ComponentName
}

func (c *SettingsComponent) WeightGet() int {
	return 50
}

/* Settings special */
func (c *SettingsComponent) RegisterDBDriver(s string, db DBSettingsApi) {
	c.drivers[s] = db
}

func (c *SettingsComponent) DBGet(e echo.Context) (DBSettingsApi, *shared.APIError) {
	if e == nil {
		return c.driver.DBGet()
	}

	return c.driver.DBFromContext(e)
}

func (c *SettingsComponent) ConfigGet() *Configuration {
	return c.config
}

func init() {
	Component = &SettingsComponent{
		drivers: make(map[string]DBSettingsApi),
	}
	registry.Instance().Register(ComponentName, Component)
}
