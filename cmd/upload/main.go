package main

import (
	"bufio"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/soltanat/otus-highload/internal/bootstrap/db"
	"github.com/soltanat/otus-highload/internal/entity"
	"github.com/soltanat/otus-highload/internal/logger"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	l := logger.Get()

	l.Info().Msg("starting upload")

	ctx := context.Background()

	parseFlags()

	var conn *pgxpool.Pool
	var err error
	for i := 0; i < 3; i++ {
		conn, err = db.New(ctx, flagDBAddr)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
			l.Fatal().Err(err).Msg("unable to connect to database")
		}
		break
	}

	c := New()
	_, err = conn.CopyFrom(ctx, pgx.Identifier{"users"}, []string{"id", "first_name", "second_name", "birthdate", "biography", "city", "password"}, c)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to copy users")
	}
}

type CopyUsers struct {
	scanner *bufio.Scanner
	file    *os.File
	ch      chan *entity.User
	count   int
}

func New() *CopyUsers {
	l := logger.Get()

	file, err := os.Open(flagFileName)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to open file")
	}

	scanner := bufio.NewScanner(file)

	return &CopyUsers{
		scanner: scanner,
		file:    file,
		ch:      make(chan *entity.User, 1),
	}
}

func (c *CopyUsers) Next() bool {
	l := logger.Get()

	for c.scanner.Scan() {
		s := c.scanner.Text()

		parts := strings.Split(s, " ")
		secondName := parts[0]

		parts = strings.Split(parts[1], ",")

		firstName := parts[0]

		age, err := strconv.Atoi(parts[1])
		if err != nil {
			l.Fatal().Err(err).Msg("unable to convert age")
		}

		birthdate := time.Now().AddDate(-age, 0, 0)

		city := parts[2]

		c.ch <- &entity.User{
			ID:         uuid.New(),
			FirstName:  &firstName,
			SecondName: &secondName,
			BirthDate:  &birthdate,
			Biography:  nil,
			City:       &city,
			Password:   nil,
		}

		return true
	}

	close(c.ch)

	return false
}

func (c *CopyUsers) Values() ([]any, error) {
	item := <-c.ch
	return []any{
		item.ID,
		item.FirstName,
		item.SecondName,
		item.BirthDate,
		item.Biography,
		item.City,
		item.Password,
	}, nil
}

func (c *CopyUsers) Err() error {
	return nil
}

func (c *CopyUsers) Close() error {
	return c.file.Close()
}
