FROM golang:1.22-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /migrations
