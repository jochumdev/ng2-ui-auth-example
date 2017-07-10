package statichttp

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/registry"
)

const (
	// ComponentName is the name of this component.
	ComponentName = "app_static_http"
)

type Configuration struct {
	Enabled   bool
	StaticDir string
	Index     string
}

type StaticHttpComponent struct {
	config *Configuration
}

// SetupIni configures lxdweb by using gopkg.in/ini.v1
func (c *StaticHttpComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	cfg := &Configuration{}

	if err := iniCfg.Section(ComponentName).MapTo(cfg); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	return c.SetupStruct(cfg, configFile, debug)
}

// SetupStruct Configures lxdweb from a struct
func (c *StaticHttpComponent) SetupStruct(cfg *Configuration, configFile string, debug bool) error {
	c.config = cfg

	if !filepath.IsAbs(c.config.StaticDir) {
		c.config.StaticDir = filepath.Join(
			filepath.Dir(configFile), c.config.StaticDir,
		)
	}

	if c.config.Index == "" {
		c.config.Index = "index.html"
	}

	return nil
}

// SetupEcho sets an echo instance up.
func (c *StaticHttpComponent) SetupEcho(e *echo.Echo) error {
	if c.config.Enabled {
		echo.NotFoundHandler = func(c2 echo.Context) error {
			index := filepath.Join(c.config.StaticDir, c.config.Index)
			_, err := os.Open(index)
			if err != nil {
				return echo.ErrNotFound
			}
			return c2.File(path.Join(c.config.StaticDir, c.config.Index))
		}

		log.Printf("Serving static content from '%s'", c.config.StaticDir)
		e.File("/", filepath.Join(c.config.StaticDir, c.config.Index))
		e.Static("/", c.config.StaticDir)
	}

	return nil
}

func (c *StaticHttpComponent) Shutdown() error {
	return nil
}

func (c *StaticHttpComponent) RegistrySet(r *registry.Registry) {}

func (c *StaticHttpComponent) NameGet() string {
	return ComponentName
}

func (c *StaticHttpComponent) WeightGet() int {
	return 1000
}

func init() {
	registry.Instance().Register(ComponentName, &StaticHttpComponent{})
}
