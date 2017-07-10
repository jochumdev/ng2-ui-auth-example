package auth

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/settings"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

type oauth2InputData struct {
	Code        string `json:"code"`
	ClientId    string `json:"clientId"`
	RedirectURI string `json:"redirectUri"`
}

func apiRegisterRoutes(e *echo.Echo) {
	ra := e.Group("/api/auth/v1")
	ra.POST("/login", apiLoginPost)

	ra.POST("/signup", apiSignupPost)

	// See: https://github.com/ronzeidman/ng2-ui-auth/issues/90#issuecomment-312504337
	ra.GET("/blank", apiBlank)

	roauth := ra.Group("/a")
	roauth.Use(auth.MiddlewareJWTAuthJSON)
	roauth.POST("/google", auth.AuthorisationWrapper(apiCallbackGooglePost, ""))
	roauth.POST("/facebook", auth.AuthorisationWrapper(apiCallbackFacebookPost, ""))
	roauth.POST("/github", auth.AuthorisationWrapper(apiCallbackGithubPost, ""))
	roauth.POST("/twitter", auth.AuthorisationWrapper(apiCallbackTwitterPost, ""))

	rme := ra.Group("/me")
	rme.Use(auth.MiddlewareJWTAuthJSON)
	rme.GET("", auth.AuthorisationWrapper(apiMe, auth.PermissionLoggedIn))
	rme.PUT("", auth.AuthorisationWrapper(apiMePut, PermissionUpdateProfile))

	rme.GET("/refresh", auth.AuthorisationWrapper(apiRefreshGet, auth.PermissionLoggedIn))
	rme.POST("/unlink", auth.AuthorisationWrapper(apiUnlinkPost, auth.PermissionLoggedIn))
}

func apiSendToken(c echo.Context, u *auth.User) error {
	token, err := u.TokenGenerate()
	if err != nil {
		return shared.APIHandleError(
			c,
			&shared.APIError{
				Reason:   err,
				Code:     http.StatusUnauthorized,
				Internal: true,
			},
		)
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})

}

func apiBlank(c echo.Context) error {
	return c.String(http.StatusOK, "Blank")
}

func apiRefreshGet(c echo.Context, u *auth.User) error {
	return apiSendToken(c, u)
}

func apiUnlinkPost(c echo.Context, u *auth.User) error {
	type unlinkData struct {
		Provider string `json:"provider"`
	}

	input := &unlinkData{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}

	if input.Provider == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Missing information"})
	}

	if _, ok := u.Properties[input.Provider]; ok {
		u.PropertyDelete(input.Provider)
	} else {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Provider not linked"})
	}

	db, apiErr := auth.DBFromContext(c)
	if apiErr != nil {
		return shared.APIHandleError(c, *apiErr)
	}
	db.UserUpdate(u)

	return apiSendToken(c, u)
}

func apiSignupPost(c echo.Context) error {

	// Check if the setting "com_settings.AllowSignup" is true
	if true {
		sdb, apiErr := settings.Component.DBGet(c)
		if apiErr != nil {
			return shared.APIHandleError(c, *apiErr)
		}

		as, apiErr := sdb.SettingGet(SettingAllowSignup)
		if apiErr != nil {
			return shared.APIHandleError(c, *apiErr)
		}
		if sv, err := as.AsBool(); !sv || err != nil {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Signup is not allowed"})
		}
	}

	type signupData struct {
		Username string `json:"username"`
		EMail    string `json:"email"`
		Password string `json:"password"`
	}

	input := &signupData{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	if input.Username == "" || input.EMail == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing information"})
	}

	db, apiErr := auth.DBFromContext(c)
	if apiErr != nil {
		return shared.APIHandleError(c, *apiErr)
	}

	user := auth.UserNew(input.Username)
	user.PasswordSet(input.Password)

	// Give the !first! local created user super admin privileges.
	roleAdded := false
	ip := c.RealIP()
	if ip == "::1" || ip == "127.0.0.1" {
		var count int
		if count, apiErr = db.UserCount(); apiErr != nil {
			return shared.APIHandleError(c, apiErr)
		}
		if count == 0 {
			roleAdded = true
			user.RoleAdd(shared.RoleSuperAdmin)
		}
	}

	// Every new user gets Guest Role.
	if !roleAdded {
		user.RoleAdd(shared.RoleGuest)
	}

	user.PropertySet("email", input.EMail)

	if apiErr := db.UserCreate(user); apiErr != nil {
		return shared.APIHandleError(c, apiErr)
	}

	return apiSendToken(c, user)
}

func apiLoginPost(c echo.Context) error {
	type userData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	input := &userData{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	if input.Username == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing information"})
	}

	db, apiErr := auth.DBFromContext(c)
	if apiErr != nil {
		return shared.APIHandleError(c, apiErr)
	}

	user, apiErr := db.UserFindByUsername(input.Username)
	if apiErr != nil {
		return shared.APIHandleError(c, apiErr)
	}

	err := user.Authenticate(input.Password)
	if err != nil {
		return shared.APIHandleError(
			c,
			&shared.APIError{
				Reason:   auth.ErrWrongUsernameOrPassword,
				Code:     http.StatusUnauthorized,
				Internal: false,
			},
		)
	}

	return apiSendToken(c, user)
}

func apiMe(c echo.Context, u *auth.User) error {
	type apiMeResult struct {
		DisplayName string `json:"displayName"`
		EMail       string `json:"email"`
		Picture     string `json:"picture"`
		LFacebook   bool   `json:"l_facebook"`
		LGoogle     bool   `json:"l_google"`
		LLinkedIn   bool   `json:"l_linkedin"`
		LTwitter    bool   `json:"l_twitter"`
		LGithub     bool   `json:"l_github"`
		LInstagram  bool   `json:"l_instagram"`
		LFoursquare bool   `json:"l_foursquare"`
		LYahoo      bool   `json:"l_yahoo"`
		LLive       bool   `json:"l_live"`
		LTwitch     bool   `json:"l_twitch"`
		LBitbucket  bool   `json:"l_bitbucket"`
		LSpotify    bool   `json:"l_spotify"`
	}

	r := &apiMeResult{
		DisplayName: u.PropertyGet("displayName", ""),
		EMail:       u.PropertyGet("email", ""),
		Picture:     u.PropertyGet("picture", ""),
		LFacebook:   u.PropertyGet("facebook", "") != "",
		LGoogle:     u.PropertyGet("google", "") != "",
		LLinkedIn:   u.PropertyGet("linkedin", "") != "",
		LTwitter:    u.PropertyGet("twitter", "") != "",
		LGithub:     u.PropertyGet("github", "") != "",
		LInstagram:  u.PropertyGet("instagram", "") != "",
		LFoursquare: u.PropertyGet("foursquare", "") != "",
		LYahoo:      u.PropertyGet("yahoo", "") != "",
		LLive:       u.PropertyGet("live", "") != "",
		LTwitch:     u.PropertyGet("twitch", "") != "",
		LBitbucket:  u.PropertyGet("bitbucket", "") != "",
		LSpotify:    u.PropertyGet("spotify", "") != "",
	}

	return c.JSON(http.StatusOK, r)
}

func apiMePut(c echo.Context, u *auth.User) error {
	type apiMeInput struct {
		DisplayName string `json:"displayName"`
		EMail       string `json:"email"`
	}

	input := &apiMeInput{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	u.PropertySet("displayName", input.DisplayName)
	u.PropertySet("email", input.EMail)

	db, err := auth.DBFromContext(c)
	if err != nil {
		return shared.APIHandleError(c, err)
	}

	if err = db.UserUpdate(u); err != nil {
		return shared.APIHandleError(c, err)
	}

	return c.JSON(http.StatusOK, "ok")
}

type apiOAuthPData struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	EMail   string `json:"email"`
	Picture string `json:"picture"`
}

func apiAfterOAuth(c echo.Context, u *auth.User, provider string, oPD *apiOAuthPData) error {
	if oPD.Id == "" {
		return c.String(http.StatusNotAcceptable, "Failed to fetch data from provider: "+provider)
	}

	db, apiEerr := auth.DBFromContext(c)
	if apiEerr != nil {
		return shared.APIHandleError(c, apiEerr)
	}

	pUser, _ := db.UserFindByProperty(provider, oPD.Id)

	// Step 3a. Link
	if u != auth.Component.UserGuest {
		if pUser != nil {
			return c.String(http.StatusConflict, provider+" profile already linked")
		}

		u.PropertySet(provider, oPD.Id)
		u.PropertySet("picture", u.PropertyGet("picture", oPD.Picture))
		u.PropertySet("displayName", u.PropertyGet("displayName", oPD.Name))
		u.PropertySet("email", u.PropertyGet("email", oPD.EMail))

		// Save the changes
		if apiEerr = db.UserUpdate(u); apiEerr != nil {
			return shared.APIHandleError(c, apiEerr)
		}
		pUser = u

	} else if pUser == nil {
		// Step 3b. Create a new user account

		// Check if signup is enabled
		sdb, apiErr := settings.Component.DBGet(c)
		if apiErr != nil {
			return shared.APIHandleError(c, apiErr)
		}

		s, apiErr := sdb.SettingGet(SettingAllowSignup)
		if apiErr != nil {
			return shared.APIHandleError(c, apiErr)
		}
		if ok, _ := s.AsBool(); !ok {
			return c.JSON(http.StatusForbidden, echo.Map{"message": "Signup is not allowed."})
		}

		pUser = auth.UserNew(provider + "|" + oPD.Id)
		pUser.RoleAdd(shared.RoleUser)
		pUser.PropertySet(provider, oPD.Id)
		pUser.PropertySet("email", oPD.EMail)
		pUser.PropertySet("picture", oPD.Picture)
		pUser.PropertySet("displayName", oPD.Name)

		// Save the new user.
		if apiEerr = db.UserUpdate(pUser); apiEerr != nil {
			return shared.APIHandleError(c, apiEerr)
		}
	}

	// Existing users are already found.

	return apiSendToken(c, pUser)
}
