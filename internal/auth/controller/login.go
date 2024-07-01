package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/gauth/internal/auth/service"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	id := c.Query("id")

	fmt.Printf("username: %q, password: %q, id: %q\n", username, password, id)
	authReq, err := service.GetAuthReqById(id)
	if err != nil {
		fmt.Printf("aerror getting auth req: %+v", err)
		c.JSON(http.StatusBadRequest, 0)
		return
	}

	fmt.Printf("req: %+v\n", authReq)
	code := authReq.GenerateCode()

	// service.
	fmt.Printf("code: %+v\n", code)
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?code=%s", authReq.RedirectUri, code))
	return
}
