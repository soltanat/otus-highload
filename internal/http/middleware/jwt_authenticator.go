package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	echoStrictMiddleware "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

type JWSValidator interface {
	ValidateJWS(tokenString string) (*jwt.Token, error)
}

const userIDKey = "subject"

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
)

func NewAuthenticator(v JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "bearerAuth" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid security scheme")
	}

	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWS")
	}

	token, err := v.ValidateJWS(jws)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWS")
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid JWS")
	}

	eCtx := middleware.GetEchoContext(ctx)
	eCtx.Set(userIDKey, sub)

	return nil
}

func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

var UserIDKeyStruct = struct{}{}

func StrictMiddlewareUserIDTransfer(f echoStrictMiddleware.StrictEchoHandlerFunc, operationID string) echoStrictMiddleware.StrictEchoHandlerFunc {
	return func(ctx echo.Context, request interface{}) (response interface{}, err error) {
		value := ctx.Get(userIDKey)
		if value != nil {
			if _, ok := value.(string); !ok {
				return nil, fmt.Errorf("user_id is not a string")
			}
			rCtx := ctx.Request().Context()
			rCtx = context.WithValue(rCtx, UserIDKeyStruct, value)
			ctx.SetRequest(ctx.Request().WithContext(rCtx))
		}

		return f(ctx, request)
	}
}
