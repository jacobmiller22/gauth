package main

import (
	"context"
	"log/slog"
	"os"

	lk "github.com/jacobmiller22/gauth/internal/logkeys"
	"github.com/jacobmiller22/gauth/pkg/clog"
)

func main() {

	l := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(lk.Loc, "main")
	ctx := clog.WithContext(context.Background(), l)
	// db := dbtest.NewDB()

	app := &gauthApp{
		ctx: ctx,
		// db:  db,
	}

	_ = app.serveHttp("0.0.0.0", 5999)
}
