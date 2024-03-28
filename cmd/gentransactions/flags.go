package main

import (
	"flag"
)

var flagWriteDBAddr string

func parseFlags() {
	flag.StringVar(&flagWriteDBAddr, "wd", "postgres://postgres:postgres@localhost:5432/social?sslmode=disable", "database dsn")
	flag.Parse()
}
