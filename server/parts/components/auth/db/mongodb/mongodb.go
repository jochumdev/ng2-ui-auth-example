package mongodb

import (
	"errors"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/mongodb"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

const driverName = "com_auth_mongodb"

type mongoDriver struct {
	db     *mgo.Database
	config *auth.Configuration
}

func (m *mongoDriver) ConfigSet(c *auth.Configuration) {
	m.config = c
}

func (m *mongoDriver) SetupEcho(e *echo.Echo) error {
	url := auth.Component.ConfigGet().AuthDBUrl
	if url == "" {
		url = mongodb.Component.ConfigGet().Url
	}

	if auth.Component.ConfigGet().AuthDBType == shared.DBTypeMongoDB {
		e.Use(mongodb.MiddlewareMongoDB(driverName, mongodb.DBConnect(url)))
	}

	return nil
}

// DBGet gets a database connection for use in commandline/testing envs.
func (m *mongoDriver) DBGet() (auth.DBAuthAPI, *shared.APIError) {
	url := m.config.AuthDBUrl
	if url == "" {
		url = mongodb.Component.ConfigGet().Url
	}

	dbSession := mongodb.DBConnect(url)
	d := &mongoDriver{db: dbSession.DB("")}
	d.ConfigSet(m.config)
	return d, nil
}

func (m *mongoDriver) DBFromContext(c echo.Context) (auth.DBAuthAPI, *shared.APIError) {
	db := c.Get(driverName).(*mgo.Database)
	return &mongoDriver{db: db}, nil
}

func (m *mongoDriver) Initialize() error {
	i := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		Background: false,
		Name:       "username",
	}
	if err := m.db.C(m.config.UsersCollection).EnsureIndex(i); err != nil {
		log.Print("ERROR: Creating Index 'username'")
	}

	i = mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		Background: false,
		Name:       "name",
	}
	if err := m.db.C(m.config.RolesCollection).EnsureIndex(i); err != nil {
		log.Print("ERROR: Creating Index 'name'")
	}

	i = mgo.Index{
		Key:        []string{"permission"},
		Unique:     true,
		Background: false,
		Name:       "permission",
	}
	if err := m.db.C(m.config.PermissionsCollection).EnsureIndex(i); err != nil {
		log.Print("ERROR: Creating Index 'permission'")
	}

	return nil
}

func (m *mongoDriver) UserFindByID(id string) (*auth.User, *shared.APIError) {
	uC := m.db.C(m.config.UsersCollection)
	mu := &auth.User{}
	err := uC.FindId(bson.ObjectIdHex(id)).One(mu)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &shared.APIError{
				Reason:   auth.ErrWrongUsernameOrPassword,
				Internal: false,
				Code:     http.StatusUnauthorized,
			}
		}
		return nil, &shared.APIError{Reason: err, Internal: true}
	}

	mu.PermissionsCacheGenerate(m)
	return mu, nil
}

func (m *mongoDriver) UserFindByUsername(username string) (*auth.User, *shared.APIError) {
	uC := m.db.C(m.config.UsersCollection)
	user := &auth.User{}
	err := uC.Find(bson.M{"username": username}).One(user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &shared.APIError{
				Reason:   auth.ErrWrongUsernameOrPassword,
				Internal: false,
				Code:     http.StatusUnauthorized,
			}
		}
		return nil, &shared.APIError{Reason: err, Internal: true}
	}

	user.PermissionsCacheGenerate(m)
	return user, nil
}

func (m *mongoDriver) UserFindByProperty(property string, value string) (*auth.User, *shared.APIError) {
	return nil, &shared.APIError{
		Reason:   errors.New("Not implemented"),
		Internal: false,
		Code:     http.StatusNotImplemented,
	}
}

func (m *mongoDriver) UserCreate(u *auth.User) *shared.APIError {
	uC := m.db.C(m.config.UsersCollection)
	err := uC.Insert(u)
	if mgo.IsDup(err) {
		return &shared.APIError{Reason: errors.New("User already exists"), Internal: false}
	}

	return nil
}

func (m *mongoDriver) UserUpdate(u *auth.User) *shared.APIError {
	uC := m.db.C(m.config.UsersCollection)

	_, err := uC.UpsertId(u.ID, bson.M{"$set": u})
	if err != nil {
		return &shared.APIError{Reason: err, Internal: true}
	}

	return nil
}

func (m *mongoDriver) UserDelete(username string) *shared.APIError {
	uC := m.db.C(m.config.UsersCollection)
	if err := uC.Remove(bson.M{"username": username}); err != nil {
		return &shared.APIError{Reason: err, Internal: true}
	}

	return nil
}

func (m *mongoDriver) UserCount() (int, *shared.APIError) {
	return 0, &shared.APIError{Reason: errors.New("Not Implemented")}
}

func (m *mongoDriver) PermissionCreate(permission string) *shared.APIError {
	pC := m.db.C(m.config.PermissionsCollection)
	if err := pC.Insert(map[string]string{"permission": permission}); err != nil {
		if mgo.IsDup(err) {
			return &shared.APIError{Reason: auth.ErrPermissionExists, Internal: false}
		}

		return &shared.APIError{Reason: err, Internal: true}
	}

	return nil
}

func (m *mongoDriver) RoleFind(name string) (*auth.Role, *shared.APIError) {
	rC := m.db.C(m.config.RolesCollection)
	r := &auth.Role{}
	err := rC.Find(bson.M{"name": name}).One(r)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &shared.APIError{
				Reason:   auth.ErrNotFound,
				Internal: false,
				Code:     http.StatusNotFound,
			}
		}
		return nil, &shared.APIError{Reason: err, Internal: true}
	}

	return r, nil
}

func (m *mongoDriver) RoleCreate(role *auth.Role) *shared.APIError {
	rC := m.db.C(m.config.RolesCollection)
	if err := rC.Insert(role); err != nil {
		if mgo.IsDup(err) {
			return &shared.APIError{Reason: auth.ErrRoleExists, Internal: false}
		}

		return &shared.APIError{Reason: err, Internal: true}
	}

	return nil
}

func (m *mongoDriver) RoleUpdate(role *auth.Role) *shared.APIError {
	rC := m.db.C(m.config.RolesCollection)
	_, err := rC.Upsert(bson.M{"name": role.Name}, role)
	if err != nil {
		return &shared.APIError{Reason: err, Internal: true}
	}

	return nil
}

func init() {
	auth.Component.RegisterDBDriver("mongodb", &mongoDriver{})
}
