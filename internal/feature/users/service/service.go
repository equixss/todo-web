package users_service

import (
	"context"

	"github.com/equixss/todo-web/internal/core/domain"
)

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

type UsersService struct {
	usersRepository UsersRepository
}

func NewUsersService(repository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: repository,
	}
}
