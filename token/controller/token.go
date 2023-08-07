package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authService "github.com/jacobmiller22/goth/auth/service"
	"github.com/jacobmiller22/goth/token/service"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func Token(c *gin.Context) {
	code := c.Query("code")
	// grantType := c.Query("fsdfgrant_type")
	clientId := c.Query("client_id")
	clientSecret := c.Query("client_secret")
	redirectUri := c.Query("redirect_uri")
	c.Header("Cache-Control", "no-store")

	authReq, err := authService.GetAuthReqByCode(code)
	if err != nil {
		fmt.Printf("error getting authreq: %+v\n", err)
		c.JSON(http.StatusBadRequest, 0)
		return
	}
	fmt.Printf("authReq: %+v\n", authReq)
	if clientId != authReq.ClientId {
		c.JSON(http.StatusBadRequest, 0)
		return
	}

	if redirectUri != authReq.RedirectUri {
		c.JSON(http.StatusBadRequest, 0)
		return
	}

	token := service.GenerateBearerToken(clientId, clientSecret)
	res := TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   time.Duration(time.Hour * 8).Milliseconds(),
	}

	c.JSON(http.StatusOK, res)
}
