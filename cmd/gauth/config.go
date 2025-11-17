package main

import (
	"flag"
	"fmt"
	"github.com/jacobmiller22/gauth/internal/version"
)

type config struct {
	http struct {
		host string
		port int
	}
	db struct {
		dsn string
	}
}

func configure() *config {

	var cfg config

	flag.StringVar(&cfg.http.host, "http-host", "localhost", "base URL for the application")
	flag.IntVar(&cfg.http.port, "http-port", 8000, "port to listen on for HTTP requests")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "db.sqlite", "sqlite3 DSN")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		panic(fmt.Sprintf("version: %s\n", version.Get()))
	}

	return &cfg
}
