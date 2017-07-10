package jsonstore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/registry"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
	ini "gopkg.in/ini.v1"
)

const (
	ComponentName = "com_jsonstore"
)

var (
	Component *JSONStoreComponent
)

// Configuration is the configuration for "auth"
type Configuration struct {
	DataDir string
}

type JSONStoreComponent struct {
	config *Configuration
	Debug  bool
}

func (c *JSONStoreComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	c.Debug = debug
	c.config = &Configuration{}
	if err := iniCfg.Section(ComponentName).MapTo(c.config); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	if c.config.DataDir == "" {
		c.config.DataDir = "json"
	}

	if !filepath.IsAbs(c.config.DataDir) {
		c.config.DataDir = filepath.Join(
			filepath.Dir(configFile), c.config.DataDir,
		)
	}

	return nil
}

func (c *JSONStoreComponent) SetupEcho(e *echo.Echo) error {
	return nil
}

func (c *JSONStoreComponent) Shutdown() error {
	return nil
}

func (c *JSONStoreComponent) RegistrySet(r *registry.Registry) {}

func (c *JSONStoreComponent) NameGet() string {
	return ComponentName
}

func (c *JSONStoreComponent) WeightGet() int {
	return 40
}

/* JSONStoreComponent exclusive */
func (c *JSONStoreComponent) ConfigGet() *Configuration {
	return c.config
}

func AbsFilePath(filename string) string {
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(Component.ConfigGet().DataDir, filename)
	}

	return filename
}

func Load(filename string, v interface{}) *shared.APIError {
	filename = AbsFilePath(filename)

	// if Component.Debug {
	// 	log.Printf("%s: Loading %s", ComponentName, filename)
	// }

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return &shared.APIError{
			Reason:   err,
			Internal: false,
			Code:     http.StatusInternalServerError,
		}
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return &shared.APIError{
			Reason:   err,
			Internal: false,
			Code:     http.StatusInternalServerError,
		}
	}

	return nil
}

func Save(filename string, v interface{}) *shared.APIError {
	filename = AbsFilePath(filename)

	// if Component.Debug {
	// 	log.Printf("%s: Saving %s", ComponentName, filename)
	// }

	data, err := json.Marshal(v)
	if err != nil {
		return &shared.APIError{
			Reason:   err,
			Internal: false,
			Code:     http.StatusInternalServerError,
		}
	}

	err = ioutil.WriteFile(filename, data, os.FileMode(0600))
	if err != nil {
		return &shared.APIError{
			Reason:   err,
			Internal: false,
			Code:     http.StatusInternalServerError,
		}
	}

	return nil
}

func init() {
	Component = &JSONStoreComponent{}
	registry.Instance().Register(ComponentName, Component)
}
