package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/ini.v1"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/parts/components/registry"
	"github.com/pcdummy/ng2-ui-auth-example/shared"
)

const (
	ComponentName = "com_auth"
)

// Configuration is the configuration for "auth"
type Configuration struct {
	Debug bool

	AuthDBType string
	AuthDBUrl  string

	AllowSignup bool

	JWTKeyFile           string
	JWTPubKeyFile        string
	JWTExpirationSeconds int

	FacebookClientID string
	FacebookSecret   string
	GoogleClientID   string
	GoogleSecret     string
	GithubClientID   string
	GithubSecret     string
	TwitterKey       string
	TwitterSecret    string

	// FoursquareSecret  string
	// LinkedinSecret    string
	// WindowsLiveSecret string
	// YahooSecret       string

	UsersCollection       string
	RolesCollection       string
	PermissionsCollection string
}

var (
	Component *AuthComponent
	signKey   *rsa.PrivateKey
	pubKey    *rsa.PublicKey
)

type AuthComponent struct {
	config   *Configuration
	registry *registry.Registry

	driver  DBAuthAPI
	drivers map[string](DBAuthAPI)

	UserGuest *User
}

func (c *AuthComponent) NameGet() string {
	return ComponentName
}

func (c *AuthComponent) WeightGet() int {
	return 100
}

func (c *AuthComponent) SetupFromIni(iniCfg *ini.File, configFile string, debug bool) error {
	cfg := &Configuration{JWTExpirationSeconds: 345600}

	if err := iniCfg.Section(ComponentName).MapTo(cfg); err != nil {
		return fmt.Errorf("Failed to parse section '%s': %v", ComponentName, err)
	}

	return c.SetupStruct(cfg, configFile, debug)
}

func (c *AuthComponent) SetupStruct(cfg *Configuration, configFile string, debug bool) error {
	c.config = cfg
	c.config.Debug = debug

	if c.config.UsersCollection == "" {
		c.config.UsersCollection = "auth_users"
	}

	if c.config.RolesCollection == "" {
		c.config.RolesCollection = "auth_roles"
	}

	if c.config.PermissionsCollection == "" {
		c.config.PermissionsCollection = "auth_permissions"
	}

	MiddlewareJWTAuth = middlewareJWTAuth(false)
	MiddlewareJWTAuthJSON = middlewareJWTAuth(true)

	var (
		driver DBAuthAPI
		ok     bool
	)

	if driver, ok = c.drivers[c.config.AuthDBType]; !ok {
		log.Fatalf(
			"Unknown db backend '%s' configured for auth.",
			c.config.AuthDBType,
		)
	}

	c.driver = driver
	c.driver.ConfigSet(c.config)

	if true {
		db, _ := DBGet()

		if err := db.Initialize(); err != nil {
			log.Fatalf("Failed to initialize the db: %v", err)
		}

		// Create default roles
		roles := []string{
			shared.RoleSuperAdmin,
			shared.RoleAdmin,
			shared.RoleUser,
			shared.RoleGuest,
		}
		for _, role := range roles {
			dbRole := &Role{Name: role}
			db.RoleCreate(dbRole)
		}

		// Create permissions and append them to roles
		// if the permission doesn't exists.
		PermissionCreate(
			PermissionLoggedIn,
			shared.RoleSuperAdmin, shared.RoleAdmin, shared.RoleUser, shared.RoleViewer,
		)
		PermissionCreate(
			PermissionGuest,
			shared.RoleSuperAdmin, shared.RoleAdmin, shared.RoleUser, shared.RoleViewer, shared.RoleGuest,
		)

		// Create the guest user
		c.UserGuest = &User{
			Username: "guest",
			Roles:    []string{shared.RoleGuest},
		}

		c.UserGuest.PermissionsCacheGenerate(db)
	}

	// Get JWT key and JWT pubkey.
	if !filepath.IsAbs(c.config.JWTKeyFile) {
		c.config.JWTKeyFile = filepath.Join(
			filepath.Dir(configFile), c.config.JWTKeyFile,
		)
	}

	if !filepath.IsAbs(c.config.JWTPubKeyFile) {
		c.config.JWTPubKeyFile = filepath.Join(
			filepath.Dir(configFile), c.config.JWTPubKeyFile,
		)
	}

	signBytes, err := ioutil.ReadFile(c.config.JWTKeyFile)
	if err != nil {
		log.Fatalf("ERROR: Cannot read private keyfile: %v", c.config.JWTKeyFile)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("ERROR: Cannot parse private keyfile: %v", c.config.JWTKeyFile)
	}

	pubBytes, err := ioutil.ReadFile(c.config.JWTPubKeyFile)
	if err != nil {
		log.Fatalf("ERROR: Reading public key: %v", c.config.JWTPubKeyFile)
	}

	pubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatalf("ERROR: Cannot parse public key: %v", c.config.JWTPubKeyFile)
	}

	return nil
}

func (c *AuthComponent) SetupEcho(e *echo.Echo) error {
	return nil
}

func (c *AuthComponent) Shutdown() error {
	return nil
}

func (c *AuthComponent) RegistrySet(r *registry.Registry) {
	c.registry = r
}

/*
 * Auth Special
 */
func (c *AuthComponent) RegisterDBDriver(s string, db DBAuthAPI) {
	c.drivers[s] = db
}

func (c *AuthComponent) DriverGet() DBAuthAPI {
	return c.driver
}

func (c *AuthComponent) ConfigGet() *Configuration {
	return c.config
}

func init() {
	Component = &AuthComponent{
		drivers: make(map[string]DBAuthAPI),
	}
	registry.Instance().Register(ComponentName, Component)
}
