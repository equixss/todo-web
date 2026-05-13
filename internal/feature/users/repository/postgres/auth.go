package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/equixss/todo-web/internal/core/domain"
	core_errors "github.com/equixss/todo-web/internal/core/errors"
	core_postgres_pool "github.com/equixss/todo-web/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, phone, email, COALESCE(password_hash, '') as password_hash
	FROM todoapp.users
	WHERE email = $1;
`
	row := r.pool.QueryRow(ctx, query, email)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Name,
		&userModel.Phone,
		&userModel.Email,
		&userModel.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with email='%s' not found:%w",
				email,
				core_errors.ErrNotFound,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Name,
		userModel.Phone,
		userModel.Email,
		userModel.PasswordHash,
	), nil
}

func (r *UsersRepository) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, name, phone, email, COALESCE(password_hash, '') as password_hash
	FROM todoapp.users
	WHERE phone = $1;
`
	row := r.pool.QueryRow(ctx, query, phone)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Name,
		&userModel.Phone,
		&userModel.Email,
		&userModel.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with phone='%s' not found:%w",
				phone,
				core_errors.ErrNotFound,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Name,
		userModel.Phone,
		userModel.Email,
		userModel.PasswordHash,
	), nil
}
