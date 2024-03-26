package main

import (
	"flag"
)

var flagDBAddr string
var flagFileName string

func parseFlags() {
	flag.StringVar(&flagDBAddr, "d", "postgres://postgres:postgres@postgres:5432/social?sslmode=disable", "database dsn")
	flag.StringVar(&flagFileName, "f", "./people.csv", "file to upload")

	flag.Parse()
}
