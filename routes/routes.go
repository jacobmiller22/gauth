package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/gauth/controllers"
	"github.com/jacobmiller22/gauth/controllers/authenticate"
)

func AddRoutes(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("templates/*")

	r = loginRoutes(r)

	r.GET("/.well-known/oauth-authorization-server", controllers.HandleConfig)

	r.GET("/oauth2/authorize", controllers.HandleAuthorize)

	r.POST("/oauth2/token", controllers.HandleToken)

	// r.GET("/oauth2/introspect", controllers.ConfigRoute)

	// r.GET("/oauth2/profile", controllers.ConfigRoute)

	return r
}

func loginRoutes(r *gin.Engine) *gin.Engine {
	r.GET("/authenticate", authenticate.LoginPage)

	r.POST("/authenticate", authenticate.LoginRes)

	return r
}
