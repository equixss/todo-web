package users_service

import (
	"context"
	"fmt"
	"time"

	"github.com/equixss/todo-web/internal/core/domain"
	core_http_middleware "github.com/equixss/todo-web/internal/core/transport/http/middleware"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UsersService) Login(
	ctx context.Context,
	credentials domain.LoginCredentials,
) (domain.LoginResult, error) {
	if err := credentials.Validate(); err != nil {
		return domain.LoginResult{}, fmt.Errorf("invalid credentials: %w", err)
	}

	var user domain.User
	var err error

	if credentials.IsEmail() {
		user, err = s.usersRepository.GetUserByEmail(ctx, credentials.Identifier)
	} else {
		user, err = s.usersRepository.GetUserByPhone(ctx, credentials.Identifier)
	}

	if err != nil {
		return domain.LoginResult{}, fmt.Errorf("user lookup: %w", err)
	}

	if !VerifyPassword(credentials.Password, user.PasswordHash) {
		return domain.LoginResult{}, fmt.Errorf("invalid password: %w", fmt.Errorf("invalid password"))
	}

	tokens, err := s.generateToken(user.ID)
	if err != nil {
		return domain.LoginResult{}, fmt.Errorf("generate token: %w", err)
	}

	return domain.NewLoginResult(user, tokens), nil
}

func (s *UsersService) generateToken(userID int) (domain.AuthTokens, error) {
	expiry, err := time.ParseDuration(s.authConfig.GetExpiry())
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("parse expiry: %w", err)
	}

	expiresAt := time.Now().Add(expiry)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.authConfig.GetIssuer(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.authConfig.GetSecret()))
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("sign token: %w", err)
	}

	return domain.NewAuthTokens(tokenString, expiresAt.Unix()), nil
}

func (s *UsersService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (domain.RefreshTokenResult, error) {
	if refreshToken == "" {
		return domain.RefreshTokenResult{}, fmt.Errorf("refresh token is required: %w", fmt.Errorf("invalid request"))
	}

	claims := &core_http_middleware.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.authConfig.GetSecret()), nil
	})

	if err != nil || !token.Valid {
		return domain.RefreshTokenResult{}, fmt.Errorf("invalid refresh token: %w", fmt.Errorf("invalid token"))
	}

	if claims.UserID <= 0 {
		return domain.RefreshTokenResult{}, fmt.Errorf("invalid user id in token: %w", fmt.Errorf("invalid token"))
	}

	user, err := s.usersRepository.GetUser(ctx, claims.UserID)
	if err != nil {
		return domain.RefreshTokenResult{}, fmt.Errorf("user not found: %w", err)
	}

	if user.ID == 0 {
		return domain.RefreshTokenResult{}, fmt.Errorf("user not found: %w", fmt.Errorf("invalid token"))
	}

	tokens, err := s.generateToken(user.ID)
	if err != nil {
		return domain.RefreshTokenResult{}, fmt.Errorf("generate token: %w", err)
	}

	return domain.NewRefreshTokenResult(tokens), nil
}
