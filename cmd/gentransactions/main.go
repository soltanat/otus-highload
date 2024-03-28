package main

import (
	"context"
	"fmt"
	"github.com/soltanat/otus-highload/internal/bootstrap/db"
	"github.com/soltanat/otus-highload/internal/logger"
	"math/rand"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	l := logger.Get()

	ctx, cancel := context.WithCancel(context.Background())

	parseFlags()

	conn, err := db.New(ctx, flagWriteDBAddr)
	if err != nil {
		l.Fatal().Msg(err.Error())
	}

	conn.Exec(ctx, "create table if not exists test_transactions (id serial primary key, user_id bigint, amount bigint)")

	var transactionsCount atomic.Int64
	var transactionsCountError atomic.Int64

	goCount := 300
	for i := 0; i < goCount; i++ {
		go func() {
			for {
				_, err := conn.Exec(ctx, "insert into test_transactions (user_id, amount) values ($1, $2)", rand.Int(), rand.Int())
				if err != nil {
					l.Err(err).Msg("unable to insert transaction")
					transactionsCountError.Add(1)
					break
				}
				transactionsCount.Add(1)
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}
		}()
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	fmt.Println(transactionsCount)
	fmt.Println(transactionsCountError)

	cancel()

}
