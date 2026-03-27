package users_postgres_repository

import "github.com/equixss/todo-web/internal/core/domain"

type UserModel struct {
	ID      int
	Version int

	Name  string
	Phone *string
}

func userDomainsFromModels(users []UserModel) []domain.User {
	usersDomain := make([]domain.User, len(users))
	for i, user := range users {
		model := domain.NewUser(user.ID, user.Version, user.Name, user.Phone)
		usersDomain[i] = model
	}
	return usersDomain
}
