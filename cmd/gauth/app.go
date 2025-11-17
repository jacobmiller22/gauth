package main

import (
	"context"
	"net/http"

	"github.com/jacobmiller22/gauth/internal/authorization"
)

type gauthRoutes struct {
	authorization *authorization.Routes
}

type gauthApp struct {
	ctx    context.Context // Represents a top level context for the gauth app instance
	r      gauthRoutes
	server *http.Server
}
