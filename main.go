package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/goth/auth"
	authService "github.com/jacobmiller22/goth/auth/service"
	"github.com/jacobmiller22/goth/token"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	authService.MakeStore()

	auth.AddRoutes(r)
	token.AddRoutes(r)
	r.Run()
}
