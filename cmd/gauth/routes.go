package main

import (
	"net/http"

	"github.com/jacobmiller22/gauth/internal/logging"
	"github.com/jacobmiller22/gauth/pkg/clog"
)

const (
	RouteAuthorize string = "GET /authorize"
	RouteToken     string = "POST /oauth2/token"
)

func (app *gauthApp) routes() http.Handler {

	mux := http.NewServeMux()

	l := clog.FromContext(app.ctx)

	mux.HandleFunc(RouteAuthorize, app.r.authorization.Authorize())
	mux.HandleFunc(RouteToken, app.r.authorization.Token())

	// mux.HandleFunc("/api/apps", app.applications)
	//
	// mux.HandleFunc("/api/apps/{clientId}", app.application)

	return logging.LogReq(l, mux)
}
