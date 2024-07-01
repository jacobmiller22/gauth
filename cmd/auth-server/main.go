package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobmiller22/gauth/internal/auth"
	authService "github.com/jacobmiller22/gauth/internal/auth/service"
	"github.com/jacobmiller22/gauth/internal/token"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	authService.MakeStore()

	auth.AddRoutes(r)
	token.AddRoutes(r)
	r.Run()
}
