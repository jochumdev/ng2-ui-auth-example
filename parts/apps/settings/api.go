package settings

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/parts/components/auth"
	"github.com/pcdummy/ng2-ui-auth-example/parts/components/settings"
	"github.com/pcdummy/ng2-ui-auth-example/shared"
)

type apiProperty struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	Edit  bool        `json:"edit"`
}

func apiRegisterRoutes(e *echo.Echo) {
	ra := e.Group("/api/settings/v1")
	ra.Use(auth.MiddlewareJWTAuthJSON)
	ra.GET("/", auth.AuthorisationWrapper(apiSettingsGet, auth.PermissionGuest))
	ra.GET("/:name", auth.AuthorisationWrapper(apiSettingGet, auth.PermissionGuest))
	ra.POST("/:name", auth.AuthorisationWrapper(apiSettingPost, PermissionUpdateSetting))
}

func apiSettingsGet(c echo.Context, u *auth.User) error {
	result := []*apiProperty{}

	sdb, apiErr := settings.Component.DBGet(c)
	if apiErr != nil {
		return shared.APIHandleError(c, *apiErr)
	}

	l, apiErr := sdb.SettingsList()
	if apiErr != nil {
		return shared.APIHandleError(c, *apiErr)
	}

	for _, setting := range l {
		if !u.PermissionHas(setting.PermissionView) {
			continue
		}

		result = append(result, &apiProperty{
			Key:   setting.Key,
			Value: setting.Value,
			Type:  setting.Type,
			Edit:  u.PermissionHas(setting.PermissionEdit),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func apiSettingGet(c echo.Context, u *auth.User) error {
	sdb, apiErr := settings.Component.DBGet(c)
	if apiErr != nil {
		return shared.APIHandleError(c, apiErr)
	}

	setting, apiErr := sdb.SettingGet(c.Param("name"))
	if apiErr != nil {
		return shared.APIHandleError(c, apiErr)
	}

	if !u.PermissionHas(setting.PermissionView) {
		// We return 404 here as we don't want users to check
		// for existing settings.
		return shared.APIHandleError(c, &shared.APIError{
			Reason:   http.ErrMissingFile,
			Internal: false,
			Code:     http.StatusNotFound,
		})
	}

	result := &apiProperty{
		Key:   setting.Key,
		Value: setting.Value,
		Type:  setting.Type,
	}

	return c.JSON(http.StatusOK, result)
}

func apiSettingPost(c echo.Context, u *auth.User) error {
	return nil
}
