package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/goth/auth/controller"
)

func AddRoutes(r *gin.Engine) {
	r.GET("/oauth2/authorize", controller.Authorize)

	r.GET("/login", controller.LoginPage)

	r.POST("/login", controller.Login)
}
