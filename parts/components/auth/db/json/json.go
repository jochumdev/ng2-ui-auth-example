package json

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/parts/components/auth"
	"github.com/pcdummy/ng2-ui-auth-example/parts/components/jsonstore"
	"github.com/pcdummy/ng2-ui-auth-example/shared"
)

type jsonData struct {
	Users       map[string]*auth.User `json:"users"`
	Roles       map[string]*auth.Role `json:"roles"`
	Permissions []string
}

type jsonDriver struct {
	config   *auth.Configuration
	filepath string

	users       map[string]*auth.User
	roles       map[string]*auth.Role
	permissions []string
}

func (m *jsonDriver) load() *shared.APIError {

	d := &jsonData{}
	if apiErr := jsonstore.Load(m.filepath, d); apiErr != nil {
		return apiErr
	}

	m.users = d.Users
	m.roles = d.Roles
	m.permissions = d.Permissions

	return nil
}

func (m *jsonDriver) save() *shared.APIError {
	d := &jsonData{
		Users:       m.users,
		Roles:       m.roles,
		Permissions: m.permissions,
	}

	if apiErr := jsonstore.Save(m.filepath, d); apiErr != nil {
		return apiErr
	}

	return nil
}

func (m *jsonDriver) ConfigSet(c *auth.Configuration) {
	m.config = c

	m.filepath = m.config.AuthDBUrl
	if m.filepath == "" {
		m.filepath = "auth.json"
	}
}

func (m *jsonDriver) SetupEcho(e *echo.Echo) error {
	return nil
}

func (m *jsonDriver) Initialize() error {
	m.users = make(map[string]*auth.User)
	m.roles = make(map[string]*auth.Role)

	m.load()
	// if apiErr := m.load(); apiErr != nil {
	// 	return apiErr.Reason
	// }

	return nil
}

func (m *jsonDriver) DBGet() (auth.DBAuthAPI, *shared.APIError) {
	return m, nil
}
func (m *jsonDriver) DBFromContext(c echo.Context) (auth.DBAuthAPI, *shared.APIError) {
	return m, nil
}

func (m *jsonDriver) UserFindByID(id string) (*auth.User, *shared.APIError) {
	var (
		user *auth.User
	)

	found := false
	for _, user = range m.users {
		if user.ID.Hex() == id {
			found = true
			break
		}
	}

	if !found {
		return nil, &shared.APIError{
			Reason:   auth.ErrWrongUsernameOrPassword,
			Internal: false,
			Code:     http.StatusUnauthorized,
		}
	}

	// Make a copy so we don't have concurent accesses to the map.
	result := auth.UserNew(user.Username)
	auth.UserDeepCopy(user, result)

	result.PermissionsCacheGenerate(m)
	return result, nil
}

func (m *jsonDriver) UserFindByUsername(username string) (*auth.User, *shared.APIError) {
	var (
		user *auth.User
		ok   bool
	)
	if user, ok = m.users[username]; !ok {
		return nil, &shared.APIError{
			Reason:   auth.ErrWrongUsernameOrPassword,
			Internal: false,
			Code:     http.StatusUnauthorized,
		}
	}

	// Make a copy so we don't have concurent accesses to the map.
	result := auth.UserNew(user.Username)
	auth.UserDeepCopy(user, result)

	result.PermissionsCacheGenerate(m)
	return result, nil
}

func (m *jsonDriver) UserFindByProperty(property string, value string) (*auth.User, *shared.APIError) {
	var user *auth.User
	ok := false
	for _, user = range m.users {
		if value2, ok2 := user.Properties[property]; ok2 {
			if value2 == value {
				ok = true
				break
			}
		}
	}

	if !ok {
		return nil, &shared.APIError{
			Reason:   auth.ErrWrongUsernameOrPassword,
			Internal: false,
			Code:     http.StatusUnauthorized,
		}
	}

	// Make a copy so we don't have concurent accesses to the map.
	result := auth.UserNew(user.Username)
	auth.UserDeepCopy(user, result)

	result.PermissionsCacheGenerate(m)
	return result, nil
}

func (m *jsonDriver) UserCreate(user *auth.User) *shared.APIError {
	if _, ok := m.users[user.Username]; ok {
		return &shared.APIError{
			Reason:   auth.ErrUserExists,
			Internal: false,
			Code:     http.StatusConflict,
		}
	}

	m.users[user.Username] = user
	return m.save()
}

func (m *jsonDriver) UserUpdate(user *auth.User) *shared.APIError {
	m.users[user.Username] = user
	return m.save()
}

func (m *jsonDriver) UserDelete(username string) *shared.APIError {
	if _, ok := m.users[username]; !ok {
		return &shared.APIError{Reason: errors.New("User not found."), Internal: true}
	}

	delete(m.users, username)

	return m.save()
}

func (m *jsonDriver) UserCount() (int, *shared.APIError) {
	return len(m.users), nil
}

// Permissions() ([]string, *shared.APIError)
func (m *jsonDriver) PermissionCreate(permission string) *shared.APIError {
	exists := false
	for _, perm := range m.permissions {
		if perm == permission {
			exists = true
			break
		}
	}

	if !exists {
		m.permissions = append(m.permissions, permission)
		return m.save()
	}

	return nil
}

//
// Roles() ([]*auth.Role, *shared.APIError)
func (m *jsonDriver) RoleFind(name string) (*auth.Role, *shared.APIError) {
	var (
		role *auth.Role
		ok   bool
	)

	if role, ok = m.roles[name]; !ok {
		return nil, &shared.APIError{
			Reason:   auth.ErrNotFound,
			Internal: false,
			Code:     http.StatusNotFound,
		}
	}

	return role, nil
}

func (m *jsonDriver) RoleCreate(role *auth.Role) *shared.APIError {
	m.roles[role.Name] = role

	return m.save()
}

func (m *jsonDriver) RoleUpdate(role *auth.Role) *shared.APIError {
	m.roles[role.Name] = role
	return m.save()
}

// RoleDelete(*Role) *shared.APIError

func init() {
	auth.Component.RegisterDBDriver("json", &jsonDriver{})
}
