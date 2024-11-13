package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthHandler struct {
	gin.HandlerFunc
}

const HOME_IDP_GIT_OAUTH_CLIENT_ID = ""
const HOME_IDP_GIT_OAUTH_CLIENT_SECRET = ""
const randomstring = "randomstring"

var jwtSecret = []byte("18df91ad-af53-42a1-80a6-adsgasdd3005")

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
	oauthConf := &oauth2.Config{
		ClientID:     HOME_IDP_GIT_OAUTH_CLIENT_ID,
		ClientSecret: HOME_IDP_GIT_OAUTH_CLIENT_SECRET,
		// RedirectURL:  fmt.Sprintf("%s://%s:%s/github/callback", scheme, host, port),
		RedirectURL: fmt.Sprintf("http://127.0.0.1:3000/github/callback"),
		Scopes:      []string{"user:email"},
		Endpoint:    github.Endpoint,
	}

	url := oauthConf.AuthCodeURL(randomstring, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	oauthConf := &oauth2.Config{
		ClientID:     HOME_IDP_GIT_OAUTH_CLIENT_ID,
		ClientSecret: HOME_IDP_GIT_OAUTH_CLIENT_SECRET,
		// RedirectURL:  fmt.Sprintf("%s://%s:%s/github/callback", scheme, host, port),
		RedirectURL: fmt.Sprintf("http://127.0.0.1:3000/github/callback"),
		Scopes:      []string{"user:email"},
		Endpoint:    github.Endpoint,
	}

	state := c.Query("state")
	if state != randomstring {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", "randomstring", state)
		http.Error(c.Writer, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := c.Query("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Error(c.Writer, "Code exchange failed", http.StatusBadRequest)
		return
	}

	client := oauthConf.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("Failed to get user info: '%s'\n", err)
		http.Error(c.Writer, "Failed to get user info", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Failed to parse user info: '%s'\n", err)
		http.Error(c.Writer, "Failed to parse user info", http.StatusBadRequest)
		return
	}

	username := userInfo["login"].(string)
	if username != "choigonyok" {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := jwtToken.SignedString(jwtSecret)

	json.NewEncoder(c.Writer).Encode(map[string]string{"token": t})
}
