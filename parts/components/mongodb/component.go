package mongodb

import (
	"fmt"

	"github.com/labstack/echo"

	"github.com/pcdummy/ng2-ui-auth-example/parts/components/registry"
	ini "gopkg.in/ini.v1"
	mgo "gopkg.in/mgo.v2"
)

const (
	ComponentName = "com_mongodb"
)

// Configuration is the configuration for "mongodb"
type Configuration struct {
	Url string
}

type MongoDBComponent struct {
	config    *Configuration
	dbSession *mgo.Session

	registry *registry.Registry
	weight   int
}

var (
	Component *MongoDBComponent
)

func (c *MongoDBComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	cfg := &Configuration{}

	if err := iniCfg.Section(ComponentName).MapTo(cfg); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	return c.SetupStruct(cfg, configFile, debug)
}

// SetupStruct Configures auth from a struct
func (c *MongoDBComponent) SetupStruct(cfg *Configuration, configFile string, debug bool) error {
	c.config = cfg

	return nil
}

func (c *MongoDBComponent) SetupEcho(e *echo.Echo) error {
	return nil
}

func (c *MongoDBComponent) Shutdown() error {
	return nil
}

func (c *MongoDBComponent) RegistrySet(r *registry.Registry) {
	c.registry = r
}

func (c *MongoDBComponent) NameGet() string {
	return ComponentName
}

func (c *MongoDBComponent) WeightGet() int {
	return 40
}

/* MongoDB Component special */
func (c *MongoDBComponent) ConfigGet() *Configuration {
	return c.config
}

func init() {
	Component = &MongoDBComponent{}
	registry.Instance().Register(ComponentName, Component)
}
