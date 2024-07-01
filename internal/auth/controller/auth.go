package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/gauth/internal/auth/service"
)

func Authorize(c *gin.Context) {
	clientId := c.Query("client_id")
	redirectUri := c.DefaultQuery("redirect_uri", "http://localhost:3000")
	responseType := c.Query("response_type")
	scope := c.Query("scope")
	state := c.Query("state")

	authReq := service.MakeAuthReq(clientId, redirectUri, responseType, scope, state)

	err := authReq.ForcefullyValidate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	id := authReq.Save() // Store this authorization request

	c.Redirect(http.StatusFound, fmt.Sprintf("%s?id=%s&state=%s", service.GetLoginUrl(), id, state))
}
