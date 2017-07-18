package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// User stores the user data.
type User struct {
	ID          bson.ObjectId     `bson:"_id" json:"_id"`
	Username    string            `bson:"username" json:"username"`
	Password    string            `bson:"password" json:"password"`
	Properties  map[string]string `bson:"properties" json:"properties"`
	Permissions []string          `bson:"permissions" json:"permissions"`
	Roles       []string          `bson:"roles" json:"roles"`

	permissionsCache map[string]bool
}

// UserDeepCopy deepcopies a user.
func UserDeepCopy(src, dst *User) {
	dst.ID = src.ID
	dst.Username = src.Username
	dst.Password = src.Password
	for k, v := range src.Properties {
		dst.Properties[k] = v
	}

	for _, perm := range src.Permissions {
		dst.Permissions = append(dst.Permissions, perm)
	}

	for _, role := range src.Roles {
		dst.Roles = append(dst.Roles, role)
	}
}

// UserNew creates a new user.
func UserNew(username string) *User {
	return &User{
		ID:          bson.NewObjectId(),
		Username:    username,
		Properties:  make(map[string]string),
		Permissions: []string{},
		Roles:       []string{},
	}
}

// Authenticate checks the users password.
func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

// PasswordSet sets the bcrypt encrypted password.
func (u *User) PasswordSet(password string) error {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(pwHash)
	return nil
}

// TokenGenerate generates a token for the user.
func (u *User) TokenGenerate() (string, error) {
	// New web token.
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set a header and a claim
	claims["ID"] = u.ID.Hex()
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Duration(Component.ConfigGet().JWTExpirationSeconds) * time.Second).Unix()

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate a JWT token, error was: %v", err)
	}

	return tokenString, nil
}

// RoleAdd adds a role to a user.
func (u *User) RoleAdd(role string) error {
	if u.RoleHas(role) {
		return fmt.Errorf(
			"user '%s' already has the role: '%s'", u.Username, role,
		)
	}

	u.Roles = append(u.Roles, role)
	return nil
}

// RoleHas checks if the user has the role.
func (u *User) RoleHas(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}

// RoleRemove removes a role from a user.
func (u *User) RoleRemove(role string) error {
	idx := -1
	for i, r := range u.Roles {
		if r == role {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("role '%s' not found", role)
	}

	u.Roles = append(u.Roles[:idx], u.Roles[idx+1:]...)

	return nil
}

// PermissionAdd adds a permission to a user.
func (u *User) PermissionAdd(permission string) error {
	if u.PermissionHas(permission) {
		return fmt.Errorf(
			"User '%s' already has the permission: '%s'", u.Username, permission,
		)
	}

	u.Permissions = append(u.Permissions, permission)
	return nil
}

// PermissionRemove removes a permission from a user.
func (u *User) PermissionRemove(permission string) error {
	idx := -1
	for i, perm := range u.Permissions {
		if perm == permission {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("permission '%s' not found", permission)
	}

	u.Permissions = append(u.Permissions[:idx], u.Permissions[idx+1:]...)

	return nil
}

// PermissionsCacheGenerate caches all permissions from the user and his roles.
func (u *User) PermissionsCacheGenerate(db DBAuthAPI) {
	u.permissionsCache = make(map[string]bool)

	for _, perm := range u.Permissions {
		u.permissionsCache[perm] = true
	}

	for _, rolename := range u.Roles {
		role, apiErr := db.RoleFind(rolename)
		if apiErr != nil {
			continue
		}

		for _, perm := range role.Permissions {
			u.permissionsCache[perm] = true
		}
	}
}

// PermissionsCacheClear clears the permission cache.
func (u *User) PermissionsCacheClear() {
	u.permissionsCache = nil
}

// PermissionHas returns if the user has the given permission.
func (u *User) PermissionHas(permission string) bool {
	if u.permissionsCache == nil {
		return false
	}

	if _, ok := u.permissionsCache[permission]; ok {
		return true
	}

	return false
}

// PropertyGet returns the given property or the defaultValue.
func (u *User) PropertyGet(name string, defaultValue string) string {
	if u.Properties == nil {
		return defaultValue
	}

	if value, ok := u.Properties[name]; ok {
		return value
	}

	return defaultValue
}

// PropertySet sets a property.
func (u *User) PropertySet(name string, value string) {
	u.Properties[name] = value
}

// PropertyDelete delets a property.
func (u *User) PropertyDelete(name string) {
	delete(u.Properties, name)
}

// Role holds a set of permissions
type Role struct {
	Name        string   `bson:"name" json:"name"`
	Permissions []string `bson:"permissions" json:"permissions"`
}

// PermissionAdd adds a permission to a role
func (r *Role) PermissionAdd(permission string) error {
	if r.PermissionHas(permission) {
		return fmt.Errorf(
			"Role '%s' already has the permission: '%s'", r.Name, permission,
		)
	}

	r.Permissions = append(r.Permissions, permission)
	return nil
}

// PermissionRemove removes a permission from a role.
func (r *Role) PermissionRemove(permission string) error {
	idx := -1
	for i, perm := range r.Permissions {
		if perm == permission {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("permission '%s' not found", permission)
	}

	r.Permissions = append(r.Permissions[:idx], r.Permissions[idx+1:]...)

	return nil
}

// PermissionHas returns if the role has the given permission.
func (r *Role) PermissionHas(permission string) bool {
	for _, perm := range r.Permissions {
		if perm == permission {
			return true
		}
	}

	return false
}
