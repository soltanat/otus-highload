FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/social ./cmd/social

FROM alpine

COPY --from=builder /app/bin/social /social
COPY ./migrations /migrations

CMD [ "/social" ]
