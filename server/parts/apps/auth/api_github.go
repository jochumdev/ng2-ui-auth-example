package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth"
	"golang.org/x/oauth2"
)

const apiGithubProfileURL = "https://api.github.com/user"

func apiCallbackGithubPost(c echo.Context, u *auth.User) error {
	input := &oauth2InputData{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	// Step 1. Exchange authorization code for access token
	secretConfig := auth.Component.ConfigGet()
	conf := &oauth2.Config{
		ClientID:     input.ClientId,
		ClientSecret: secretConfig.GithubSecret,
		RedirectURL:  input.RedirectURI,
		Scopes:       []string{"public_profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	t, err := conf.Exchange(context.TODO(), input.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "exchange oauth token failed: "+err.Error())
	}

	// Step 2. Retrieve profile information about the current user.
	hClient := &http.Client{}
	req, _ := http.NewRequest("GET", apiGithubProfileURL, nil)
	req.Header.Set("Authorization", "Bearer "+t.AccessToken)
	res, err := hClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch the profile: "+err.Error())
	}
	defer res.Body.Close()

	// js := &githubProfile{}
	js := make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&js)

	oPD := &apiOAuthPData{
		Id:      fmt.Sprintf("%f", js["id"].(float64)),
		Name:    js["name"].(string),
		EMail:   js["email"].(string),
		Picture: js["avatar_url"].(string),
	}
	// //
	// // Step 3. Link, Create or return the current user's JWT.
	return apiAfterOAuth(c, u, "github", oPD)
}
