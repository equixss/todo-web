package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

var EmailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type LoginCredentials struct {
	Identifier string
	Password    string
}

func NewLoginCredentials(identifier, password string) LoginCredentials {
	return LoginCredentials{
		Identifier: identifier,
		Password:    password,
	}
}

func (c *LoginCredentials) IsEmail() bool {
	return EmailRE.MatchString(c.Identifier)
}

func (c *LoginCredentials) IsPhone() bool {
	return RussianPhoneRE.MatchString(c.Identifier)
}

func (c *LoginCredentials) Validate() error {
	if c.Identifier == "" {
		return fmt.Errorf("identifier is required: %w", core_errors.ErrInvalidArgument)
	}
	if c.Password == "" {
		return fmt.Errorf("password is required: %w", core_errors.ErrInvalidArgument)
	}
	if !c.IsEmail() && !c.IsPhone() {
		return fmt.Errorf("identifier must be valid email or phone: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

type AuthTokens struct {
	AccessToken string
	ExpiresAt   int64
}

func NewAuthTokens(accessToken string, expiresAt int64) AuthTokens {
	return AuthTokens{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}
}

type LoginResult struct {
	User   User
	Tokens AuthTokens
}

func NewLoginResult(user User, tokens AuthTokens) LoginResult {
	return LoginResult{
		User:   user,
		Tokens: tokens,
	}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (r *RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return fmt.Errorf("refresh token is required: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

type RefreshTokenResult struct {
	Tokens AuthTokens
}

func NewRefreshTokenResult(tokens AuthTokens) RefreshTokenResult {
	return RefreshTokenResult{Tokens: tokens}
}
