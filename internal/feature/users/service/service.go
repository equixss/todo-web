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
	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id int,
	) error
	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (domain.User, error)
	GetUserByPhone(
		ctx context.Context,
		phone string,
	) (domain.User, error)
}

type AuthConfig interface {
	GetSecret() string
	GetExpiry() string
	GetIssuer() string
}

type UsersService struct {
	usersRepository UsersRepository
	authConfig      AuthConfig
}

func NewUsersService(repository UsersRepository, authConfig AuthConfig) *UsersService {
	return &UsersService{
		usersRepository: repository,
		authConfig:      authConfig,
	}
}
