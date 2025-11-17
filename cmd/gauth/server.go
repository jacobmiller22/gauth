package main

import (
	"fmt"
	"net/http"

	lk "github.com/jacobmiller22/gauth/internal/logkeys"
	"github.com/jacobmiller22/gauth/pkg/clog"
)

func (app *gauthApp) serveHttp(host string, port int) error {
	l := clog.FromContext(app.ctx).With(lk.Loc, "gauthApp::serveHttp")

	app.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: app.routes(),
	}
	l.Info(lk.HttpServerInitStart, "addr", app.server.Addr)

	if err := app.server.ListenAndServe(); err != nil {
		l.Error(lk.HttpServerInitFailed, lk.Error, err)
		return err
	}

	return nil
}
