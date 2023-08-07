package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/gauth/auth"
	authService "github.com/jacobmiller22/gauth/auth/service"
	"github.com/jacobmiller22/gauth/token"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	authService.MakeStore()

	auth.AddRoutes(r)
	token.AddRoutes(r)
	r.Run()
}
