package main

import (
	"flag"
)

var flagDBAddr string
var flagUsersFileName string
var flagPostsFileName string

func parseFlags() {
	flag.StringVar(&flagDBAddr, "d", "postgres://postgres:postgres@postgres:5432/social?sslmode=disable", "database dsn")
	flag.StringVar(&flagUsersFileName, "uf", "./people.csv", "file users to upload")
	flag.StringVar(&flagPostsFileName, "pf", "./posts.txt", "file posts to upload")

	flag.Parse()
}
