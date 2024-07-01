package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authService "github.com/jacobmiller22/gauth/internal/auth/service"
	"github.com/jacobmiller22/gauth/internal/token/service"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TokenError int32

const (
	InvalidRequest = iota
	InvalidClient
	InvalidGrant
	InvalidScope
	UnauthorizedClient
	UnsupportedGrantType
)

func (e TokenError) String() string {
	switch e {
	case InvalidRequest:
		return "invalid_request"
	case InvalidClient:
		return "invalid_client"
	case InvalidGrant:
		return "invalid_grant"
	case InvalidScope:
		return "invalid_scope"
	case UnauthorizedClient:
		return "unauthorized_client"
	case UnsupportedGrantType:
		return "unsupported_grant_type"
	}
	return ""
}

type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
}

func Token(c *gin.Context) {
	code := c.Query("code")
	grantType := c.Query("grant_type")
	clientId := c.Query("client_id")
	clientSecret := c.Query("client_secret")
	redirectUri := c.Query("redirect_uri")
	c.Header("Cache-Control", "no-store")

	authReq, err := authService.GetAuthReqByCode(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, TokenErrorResponse{
			Error:            TokenError.String(InvalidRequest),
			ErrorDescription: "Could not find auth code.",
		})
		return
	}
	if grantType == "client_credientials" {
		c.JSON(http.StatusBadRequest, TokenErrorResponse{
			Error:            TokenError.String(InvalidGrant),
			ErrorDescription: "Invalid grant_type",
		})
		return
	}
	if clientId != authReq.ClientId {
		c.JSON(http.StatusBadRequest, TokenErrorResponse{
			Error:            TokenError.String(InvalidClient),
			ErrorDescription: "Invalid client_id",
		})
		return
	}

	if redirectUri != authReq.RedirectUri {
		c.JSON(http.StatusBadRequest, TokenErrorResponse{
			Error:            TokenError.String(InvalidClient),
			ErrorDescription: "Mismatching redirect_uri",
		})
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
