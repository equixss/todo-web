package users_postgres_repository

type UserModel struct {
	ID      int
	Version int

	Name  string
	Phone *string
}
