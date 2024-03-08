package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	oapiEchoMiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/soltanat/otus-highload/internal/bootstrap/db"
	"github.com/soltanat/otus-highload/internal/http/api"
	"github.com/soltanat/otus-highload/internal/http/handler"
	"github.com/soltanat/otus-highload/internal/http/middleware"
	"github.com/soltanat/otus-highload/internal/logger"
	"github.com/soltanat/otus-highload/internal/storage/postgres"
	"github.com/soltanat/otus-highload/internal/usecase"
)

func main() {
	l := logger.Get()

	ctx := context.Background()

	parseFlags()

	conn, err := db.New(ctx, flagDBAddr)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to connect to database")
	}

	err = db.ApplyMigrations(flagDBAddr)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to apply migrations")
	}

	userStorage := postgres.NewUserStorage(conn)

	passHasher := usecase.NewPasswordHasher()

	userUseCase, err := usecase.NewUser(userStorage, passHasher)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to create user usecase")
	}
	tokenProvider := middleware.NewJWTProvider(flagSignatureKey)

	h := handler.New(
		userUseCase,
		tokenProvider,
	)
	strictHandler := api.NewStrictHandler(h, []api.StrictMiddlewareFunc{middleware.StrictMiddlewareUserIDTransfer})

	spec, err := api.GetSwagger()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to get swagger spec")
	}

	validator := oapiEchoMiddleware.OapiRequestValidatorWithOptions(spec,
		&oapiEchoMiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: middleware.NewAuthenticator(tokenProvider),
			},
		},
	)

	e := echo.New()
	e.HideBanner = true
	e.Use(validator)
	api.RegisterHandlers(e, strictHandler)

	go func() {
		err := e.Start(flagAddr)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			l.Fatal().Err(err).Str("addr", flagAddr).Msg("unable to start server")
		}
	}()

	gracefulShutdown()

	err = e.Close()
	if err != nil {
		l.Error().Err(err).Msg("unable to close server")
	}

}

func gracefulShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
}
