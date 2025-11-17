package main

import (
	"context"
	"net/http"

	"github.com/jacobmiller22/gauth/internal/oauth2"
)

type gauthRoutes struct {
	oauth2 *oauth2.Routes
}

type gauthApp struct {
	ctx    context.Context // Represents a top level context for the gauth app instance
	r      gauthRoutes
	server *http.Server
}
