package token

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/gauth/internal/token/controller"
)

func AddRoutes(r *gin.Engine) {
	r.GET("/oauth2/token", controller.Token)
}
