package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/parts/components/auth"
	"golang.org/x/oauth2"
)

const apiFacebookProfileURL = "https://graph.facebook.com/v2.9/me?fields=id,name,email"

func apiCallbackFacebookPost(c echo.Context, u *auth.User) error {
	input := &oauth2InputData{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	// Step 1. Exchange authorization code for access token
	secretConfig := auth.Component.ConfigGet()
	conf := &oauth2.Config{
		ClientID:     input.ClientId,
		ClientSecret: secretConfig.FacebookSecret,
		RedirectURL:  input.RedirectURI,
		Scopes:       []string{"public_profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/dialog/oauth",
			TokenURL: "https://graph.facebook.com/oauth/access_token",
		},
	}

	t, err := conf.Exchange(context.TODO(), input.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "exchange oauth token failed: "+err.Error())
	}

	// Step 2. Retrieve profile information about the current user.
	hClient := &http.Client{}
	req, _ := http.NewRequest("GET", apiFacebookProfileURL, nil)
	req.Header.Set("Authorization", "Bearer "+t.AccessToken)
	res, err := hClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch the profile: "+err.Error())
	}
	defer res.Body.Close()

	oPD := &apiOAuthPData{}
	json.NewDecoder(res.Body).Decode(&oPD)
	oPD.Picture = "https://graph.facebook.com/v2.9/" + oPD.Id + "/picture?type=large"

	// Step 3. Link, Create or return the current user's JWT.
	return apiAfterOAuth(c, u, "facebook", oPD)
}
