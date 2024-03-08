package main

import (
	"flag"
)

var flagAddr string
var flagDBAddr string
var flagSignatureKey string

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "0.0.0.0:8080", "address and port of server")
	flag.StringVar(&flagDBAddr, "d", "postgres://postgres:postgres@postgres:5432/social?sslmode=disable", "database dsn")
	flag.StringVar(&flagSignatureKey, "k", "very-secret-key", "key for signature token")

	flag.Parse()
}
