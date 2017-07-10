package json

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/jsonstore"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/settings"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

type jsonDriver struct {
	config   *settings.Configuration
	filepath string

	settings map[string]*shared.Property
}

func (m *jsonDriver) load() *shared.APIError {
	if apiErr := jsonstore.Load(m.filepath, &m.settings); apiErr != nil {
		return apiErr
	}

	return nil
}

func (m *jsonDriver) save() *shared.APIError {
	if apiErr := jsonstore.Save(m.filepath, &m.settings); apiErr != nil {
		return apiErr
	}

	return nil
}

func (m *jsonDriver) ConfigSet(c *settings.Configuration) {
	m.config = c

	m.filepath = m.config.SettingsDBUrl
	if m.filepath == "" {
		m.filepath = "settings.json"
	}
}

func (m *jsonDriver) Initialize() error {
	m.settings = make(map[string]*shared.Property)

	if apiErr := m.load(); apiErr != nil {
		return apiErr.Reason
	}

	return nil
}

func (m *jsonDriver) DBGet() (settings.DBSettingsApi, *shared.APIError) {
	return m, nil
}

func (m *jsonDriver) DBFromContext(c echo.Context) (settings.DBSettingsApi, *shared.APIError) {
	return m, nil
}

func (m *jsonDriver) SettingsList() ([]*shared.Property, *shared.APIError) {
	l := []*shared.Property{}
	for _, s := range m.settings {
		l = append(l, s)
	}

	return l, nil
}

func (m *jsonDriver) SettingsMap() (map[string]*shared.Property, *shared.APIError) {
	return m.settings, nil
}

func (m *jsonDriver) SettingGet(name string) (*shared.Property, *shared.APIError) {
	var (
		setting *shared.Property
		ok      bool
	)

	if setting, ok = m.settings[name]; !ok {
		return nil, &shared.APIError{
			Reason:   errors.New("Not found"),
			Internal: false,
			Code:     http.StatusNotFound,
		}
	}

	return setting, nil
}

func (m *jsonDriver) SettingHas(s *shared.Property) (bool, *shared.APIError) {
	_, ok := m.settings[s.Key]

	return ok, nil
}

func (m *jsonDriver) SettingCreate(s *shared.Property, overwrite bool) *shared.APIError {
	exists, apiErr := m.SettingHas(s)
	if apiErr != nil {
		return apiErr
	}
	if exists && !overwrite {
		return nil
	}

	m.settings[s.Key] = s
	return m.save()
}

func (m *jsonDriver) SettingUpdate(s *shared.Property) *shared.APIError {
	if ok, _ := m.SettingHas(s); !ok {
		return &shared.APIError{
			Reason:   errors.New("Setting not found"),
			Internal: false,
			Code:     http.StatusInternalServerError,
		}
	}

	m.settings[s.Key] = s
	return nil
}

func init() {
	settings.Component.RegisterDBDriver("json", &jsonDriver{})
}
