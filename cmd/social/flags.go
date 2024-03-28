package main

import (
	"flag"
)

var flagAddr string
var flagWriteDBAddr string
var flagReadDBAddr string
var flagSignatureKey string

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "0.0.0.0:8080", "address and port of server")
	flag.StringVar(&flagWriteDBAddr, "wd", "postgres://postgres:postgres@postgres:5432/social?sslmode=disable", "database dsn")
	//flag.StringVar(&flagReadDBAddr, "rd", "postgres://postgres:postgres@postgres-replica-0:5432,postgres-replica-1:5432/social?sslmode=disable", "database dsn")
	flag.StringVar(&flagReadDBAddr, "rd", "postgres://postgres:postgres@pgbouncer:5432/social?sslmode=disable", "database dsn")
	flag.StringVar(&flagSignatureKey, "k", "very-secret-key", "key for signature token")

	flag.Parse()
}
