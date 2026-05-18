package core_http_middleware

import (
	"context"
	"fmt"
	"strings"

	core_auth "github.com/equixss/todo-web/internal/core/auth"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_http_response "github.com/equixss/todo-web/internal/core/transport/http/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

func UserIDToContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

type JWTMiddleware struct {
	config    core_auth.Config
	presenter *core_http_response.HTTPResponsePresenter
}

func NewJWTMiddleware(config core_auth.Config, presenter *core_http_response.HTTPResponsePresenter) *JWTMiddleware {
	return &JWTMiddleware{config: config, presenter: presenter}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (m *JWTMiddleware) Authenticate() Middleware {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.presenter.ErrorResponse(
				c,
				fmt.Errorf("missing authorization header: %w", core_errors.ErrUnauthorized),
				"authorization required",
			)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			m.presenter.ErrorResponse(
				c,
				fmt.Errorf("invalid authorization header format: %w", core_errors.ErrUnauthorized),
				"authorization required",
			)
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.config.Secret), nil
		})

		if err != nil || !token.Valid {
			m.presenter.ErrorResponse(
				c,
				fmt.Errorf("invalid token: %w", core_errors.ErrUnauthorized),
				"authorization required",
			)
			return
		}

		if claims.UserID <= 0 {
			m.presenter.ErrorResponse(
				c,
				fmt.Errorf("invalid user id in token: %w", core_errors.ErrUnauthorized),
				"authorization required",
			)
			return
		}

		ctx := UserIDToContext(c.Request.Context(), claims.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
