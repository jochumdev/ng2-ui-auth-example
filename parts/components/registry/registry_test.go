package registry

import (
	"testing"

	"gopkg.in/ini.v1"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestComponent struct {
	weight int
}

func (c *TestComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	return nil
}
func (c *TestComponent) SetupEcho(e *echo.Echo) error {
	return nil
}

func (c *TestComponent) Shutdown() error {
	return nil
}

func (c *TestComponent) SettingsGet() []*shared.Property {
	return []*shared.Property{}
}

func (c *TestComponent) RegistrySet(*Registry) {}

func (c *TestComponent) NameGet() string {
	return "test"
}

func (c *TestComponent) WeightGet() int {
	return c.weight
}

type SuiteRegistry struct {
	suite.Suite
}

func (suite *SuiteRegistry) TestRegister() {
	t1 := &TestComponent{weight: 1}
	t3 := &TestComponent{weight: 3}
	t76 := &TestComponent{weight: 76}
	t75 := &TestComponent{weight: 75}
	t80 := &TestComponent{weight: 80}
	t100 := &TestComponent{weight: 100}

	r := Instance()
	r.Register("t100", t100)
	r.Register("t1", t1)
	r.Register("t80", t80)
	r.Register("t75", t75)
	r.Register("t3", t3)
	r.Register("t76", t76)

	should := []Component{t1, t3, t75, t76, t80, t100}

	assert.Equal(suite.T(), r.List(), should)
}

func (suite *SuiteRegistry) TearDownSuite() {
	instance = nil
}

func TestRegistrySuite(t *testing.T) {
	suite.Run(t, new(SuiteRegistry))
}
