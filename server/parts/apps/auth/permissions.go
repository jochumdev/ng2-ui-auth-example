package auth

const (
	PermissionUserCreate    = "app_auth.UserCreate"
	PermissionUserDelete    = "app_auth.UserDelete"
	PermissionUpdateProfile = "app_auth.UpdateProfile"

	PermissionEditSettings = "app_auth.EditSettings"
	// This permission never gets registered so its not assignable to any
	// Roles or users.
	PermissionReadonlySetting = "app_auth.ReadonlySetting"
)
