package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/gauth/service"
)

type OAuthConfiguration struct {
	SupportedScopes       []string `json:"supported_scopes"`
	Issuer                string   `json:"issuer"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	TokenEndpoint         string   `json:"token_endpoint"`
	UserInfoEndpoint      string   `json:"user_info_endpoint"`
	RevocationEndpoint    string   `json:"revocation_endpoint"`
}

func HandleAuthenticate(c *gin.Context) {
	c.HTML(http.StatusOK, "authenticate.html", gin.H{
		"title": "My title",
	})
}

func HandleConfig(c *gin.Context) {
	fp, err := os.Open("oauth.json")
	if err != nil {
		c.Status(500)
		return
	}
	fi, err := fp.Stat() // Stat gives a FileInfo object
	c.DataFromReader(http.StatusOK, fi.Size(), "application/json", fp, map[string]string{})
	return
}

func HandleToken(c *gin.Context) {
	// Code := c.Query("code")
	// GrantType := c.Query("grant_type")
	ClientId := c.Query("client_id")
	ClientSecret := c.Query("client_secret")

	token := service.GenerateBearerToken(ClientId, ClientSecret)
	res := TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   0,
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, res)
}
