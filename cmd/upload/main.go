package main

import (
	"bufio"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/soltanat/otus-highload/internal/bootstrap/db"
	"github.com/soltanat/otus-highload/internal/logger"
	"math/rand"
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

	//Users
	usersCh := make(chan []any)
	go genUsers(flagUsersFileName, usersCh)

	copyUsers := pgx.CopyFromFunc(func() (row []any, err error) {
		row = <-usersCh
		if row == nil {
			return nil, nil
		}
		return row, nil
	})
	_, err = conn.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"id", "first_name", "second_name", "birthdate", "biography", "city", "password"}, copyUsers,
	)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to copy users")
	}

	//Friends
	rows, err := conn.Query(ctx, "select users.id from users")
	if err != nil {
		l.Fatal().Err(err).Msg("unable to select users")
	}
	defer rows.Close()

	usersIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			l.Fatal().Err(err).Msg("unable to scan users")
		}
		usersIDs = append(usersIDs, userID)
	}

	friendsCh := make(chan []any)
	go genFriends(usersIDs, friendsCh)

	copyFriends := pgx.CopyFromFunc(func() (row []any, err error) {
		row = <-friendsCh
		if row == nil {
			return nil, nil
		}
		return row, nil
	})

	_, err = conn.CopyFrom(ctx, pgx.Identifier{"friends"}, []string{"user_id", "friend_id"}, copyFriends)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to copy friends")
	}

	//Stars
	rows, err = conn.Query(ctx, "select user_id from friends group by user_id having count(user_id) > 1000;")
	if err != nil {
		l.Fatal().Err(err).Msg("unable to select friends")
	}
	defer rows.Close()

	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			l.Fatal().Err(err).Msg("unable to scan users")
		}
		_, err = conn.Exec(ctx, "update users set star = true where id = $1", userID)
		if err != nil {
			l.Fatal().Err(err).Msg("unable to update friends")
		}

		_, err = conn.Exec(ctx, "update friends set star = true where user_id = $1", userID)
		if err != nil {
			l.Fatal().Err(err).Msg("unable to update friends")
		}

		l.Info().Fields(map[string]interface{}{"star_id": userID}).Msg("updated star")
	}

	l.Info().Fields(map[string]interface{}{"user_id": usersIDs[0]}).Msg("first user")

	//Posts
	rows, err = conn.Query(ctx, "select distinct friend_id from friends")
	if err != nil {
		l.Fatal().Err(err).Msg("unable to select friends")
	}
	defer rows.Close()

	friendsIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			l.Fatal().Err(err).Msg("unable to scan users")
		}
		friendsIDs = append(friendsIDs, userID)
	}

	postsCh := make(chan []any)
	go genPosts(flagPostsFileName, friendsIDs, postsCh)

	copyPosts := pgx.CopyFromFunc(func() (row []any, err error) {
		row = <-postsCh
		if row == nil {
			return nil, nil
		}
		return row, nil
	})

	_, err = conn.CopyFrom(ctx, pgx.Identifier{"posts"}, []string{"id", "author_id", "created_at", "text"}, copyPosts)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to copy posts")
	}

}

func genUsers(filePath string, ch chan []any) {
	l := logger.Get()

	file, err := os.Open(filePath)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to open file")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()

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

		ch <- []any{
			uuid.New(),
			firstName,
			secondName,
			birthdate,
			nil,
			city,
			"$2a$10$b7bEoH9IoUI3dZfe0IH0q.PCzF1YSNVAEbrsHDBgYUi6ikX5t2IR.",
		}
	}

	close(ch)
}

func genFriends(usersIDs []uuid.UUID, ch chan []any) {
	l := logger.Get()

	start, end := 0, 0
	for i := 0; i < 1000; i++ {
		userID := usersIDs[i]

		l.Info().Fields(map[string]interface{}{"count": i}).Msg("finished")

		friendsCount := rand.Intn(100)

		if i == 1 {
			friendsCount = 1000
		}

		if friendsCount == 0 {
			friendsCount = 1
		}

		start = 0
		end = rand.Intn(len(usersIDs) - 1)
		if end < friendsCount {
			end = friendsCount
		}
		if end > friendsCount {
			start = end - friendsCount
		}

		friends := usersIDs[start:end]

		if i == 1 {
			friends = append(friends, usersIDs[0])
		}

		for _, friend := range friends {
			ch <- []any{userID, friend}
		}
	}

	close(ch)
}

func genPosts(filePath string, usersIDs []uuid.UUID, ch chan []any) {
	textCh := make(chan string)
	done := make(chan struct{})
	go genPostsText(filePath, textCh, done)

	for _, userID := range usersIDs {
		postsCount := rand.Intn(50)
		for i := 0; i < postsCount; i++ {
			ch <- []any{
				uuid.New(),
				userID,
				time.Now(),
				<-textCh,
			}
		}
	}

	done <- struct{}{}
	close(textCh)
	close(done)
	close(ch)
}

func genPostsText(filePath string, ch chan string, done chan struct{}) {
	l := logger.Get()

	for {
		file, err := os.Open(filePath)
		if err != nil {
			l.Fatal().Err(err).Msg("unable to open file")
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s := scanner.Text()
			select {
			case <-done:
				file.Close()
				return
			case ch <- s:
			}
		}
	}
}
