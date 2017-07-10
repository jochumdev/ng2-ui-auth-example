package auth

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mrjones/oauth"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth"
)

const apiTwitterProfileUrl = "https://api.twitter.com/1.1/account/verify_credentials.json"

func apiCallbackTwitterPost(c echo.Context, u *auth.User) error {
	authConfig := auth.Component.ConfigGet()

	oc := oauth.NewConsumer(
		authConfig.TwitterKey,
		authConfig.TwitterSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	var requestPayload struct {
		OAuthToken    string `json:"oauth_token"`
		OAuthVerifier string `json:"oauth_verifier"`
	}
	if err := c.Bind(&requestPayload); err != nil && err != io.EOF {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}

	// Part 1/2: Initial request from Satellizer.
	if requestPayload.OAuthToken == "" || requestPayload.OAuthVerifier == "" {

		// Step 1. Obtain request token for the authorization popup.
		requestToken, _, err := oc.GetRequestTokenAndUrl(c.QueryString())
		if err != nil {
			log.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"oauth_token":              requestToken.Token,
			"oauth_token_secret":       requestToken.Secret,
			"oauth_callback_confirmed": "true",
		})
	}

	// Part 2/2: Second request after Authorize app is clicked.
	requestToken := &oauth.RequestToken{
		Token:  requestPayload.OAuthToken,
		Secret: authConfig.TwitterSecret,
	}

	// Step 3. Exchange oauth token and oauth verifier for access token.
	accessToken, err := oc.AuthorizeToken(requestToken, requestPayload.OAuthVerifier)
	if err != nil {
		log.Println(err)
		return err
	}

	// Step 4. Retrieve profile information about the current user.
	response, err := oc.Get(
		apiTwitterProfileUrl,
		map[string]string{"skip_status": "true", "include_entities": "false", "include_email": "true"},
		accessToken)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	js := make(map[string]interface{})
	json.NewDecoder(response.Body).Decode(&js)

	oPD := &apiOAuthPData{
		Id:      js["screen_name"].(string),
		Name:    js["name"].(string),
		EMail:   js["email"].(string),
		Picture: js["profile_image_url_https"].(string),
	}

	// Step 5. Link, Create or return the current user's JWT.
	return apiAfterOAuth(c, u, "twitter", oPD)
}
