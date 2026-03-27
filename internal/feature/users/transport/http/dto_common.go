package users_transport_http

import "github.com/equixss/todo-web/internal/core/domain"

type UserDTOResponse struct {
	ID      int     `json:"id"`
	Version int     `json:"version"`
	Name    string  `json:"name"`
	Phone   *string `json:"phone"`
}

func UserDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		Name:    user.Name,
		Phone:   user.Phone,
		ID:      user.ID,
		Version: user.Version,
	}
}

func UsersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))
	for i, user := range users {
		usersDTO[i] = UserDTOFromDomain(user)
	}
	return usersDTO
}
