FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o bin/upload ./cmd/upload

RUN wget https://raw.githubusercontent.com/OtusTeam/highload/master/homework/people.csv

FROM alpine

COPY --from=builder /app/bin/upload /upload
COPY --from=builder /app/people.csv /people.csv

CMD [ "/upload" ]
