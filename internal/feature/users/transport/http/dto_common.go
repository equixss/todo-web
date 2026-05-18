package users_transport_http

import "github.com/equixss/todo-web/internal/core/domain"

type UserDTOResponse struct {
	ID      int     `json:"id"`
	Version int     `json:"version"`
	Name    string  `json:"name"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email,omitempty"`
}

func UserDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:      user.ID,
		Version: user.Version,
		Name:    user.Name,
		Phone:   user.Phone,
		Email:   user.Email,
	}
}

func UsersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))
	for i, user := range users {
		usersDTO[i] = UserDTOFromDomain(user)
	}
	return usersDTO
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User        UserDTOResponse `json:"user"`
	AccessToken string          `json:"access_token"`
	ExpiresAt   int64           `json:"expires_at"`
}

func LoginResponseFromDomain(result domain.LoginResult) LoginResponse {
	return LoginResponse{
		User:        UserDTOFromDomain(result.User),
		AccessToken: result.Tokens.AccessToken,
		ExpiresAt:   result.Tokens.ExpiresAt,
	}
}

type AuthTokensDTO struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func AuthTokensFromDomain(tokens domain.AuthTokens) AuthTokensDTO {
	return AuthTokensDTO{
		AccessToken: tokens.AccessToken,
		ExpiresAt:   tokens.ExpiresAt,
	}
}
