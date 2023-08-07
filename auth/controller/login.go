package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/goth/auth/service"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	id := c.PostForm("id")

	fmt.Printf("username: %q, password: %q, id: %q", username, password, id)
	authReq, err := service.GetAuthReqById(id)
	if err != nil {
		fmt.Printf("aerror getting auth req: %+v", err)
		c.JSON(http.StatusBadRequest, 0)
		return
	}

	code := authReq.GenerateCode()

	// service.

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?code=%s", authReq.RedirectUri, code))
}
