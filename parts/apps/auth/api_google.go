package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/parts/components/auth"
	"golang.org/x/oauth2"
)

func apiCallbackGooglePost(c echo.Context, u *auth.User) error {

	input := &oauth2InputData{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	// Step 1. Exchange authorization code for access token
	secretConfig := auth.Component.ConfigGet()
	conf := &oauth2.Config{
		ClientID:     input.ClientId,
		ClientSecret: secretConfig.GoogleSecret,
		RedirectURL:  input.RedirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	t, err := conf.Exchange(context.TODO(), input.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "exchange oauth token failed: "+err.Error())
	}

	// Step 2. Retrieve profile information about the current user.
	hClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/plus/v1/people/me/openIdConnect", nil)
	req.Header.Set("Authorization", "Bearer "+t.AccessToken)
	res, err := hClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch the profile: "+err.Error())
	}
	defer res.Body.Close()

	type googlePlusProfile struct {
		Kind          string `json:"kind"`
		Gender        string `json:"gender"`
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Profile       string `json:"profile"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Locale        string `json:"locale"`
		Hd            string `json:"hd"`
	}

	gPP := &googlePlusProfile{}
	json.NewDecoder(res.Body).Decode(&gPP)

	// Step 3. Link, Create or return the current user's JWT.
	oPD := &apiOAuthPData{
		Id:      gPP.Sub,
		Name:    gPP.Name,
		EMail:   gPP.Email,
		Picture: gPP.Picture,
	}

	return apiAfterOAuth(c, u, "google", oPD)
}
