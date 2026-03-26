package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/equixss/todo-web/internal/core/errors"
)

var RussianPhoneRE = regexp.MustCompile(
	`^(?:\+7|8)\s*(?:\(?\d{3}\)?[\s\-.]?\d{3}[\s\-.]?\d{2}[\s\-.]?\d{2})$`,
)

type User struct {
	ID      int
	Version int

	Name  string
	Phone *string
}

func NewUser(
	id int,
	version int,
	name string,
	phone *string,
) User {
	return User{Name: name, Phone: phone, Version: version, ID: id}
}

func (u *User) Validate() error {
	nameLength := len([]rune(u.Name))
	if nameLength < 3 || nameLength > 100 {
		return fmt.Errorf("invalid name length: %d: %w", nameLength, core_errors.ErrInvalidArgument)
	}

	if u.Phone != nil {
		phoneLen := len([]rune(*u.Phone))
		if phoneLen < 10 || !RussianPhoneRE.MatchString(*u.Phone) {
			return fmt.Errorf("invalid phone number: %s: %w", *u.Phone, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

func NewUserUninitialized(name string, phone *string) User {
	return NewUser(UninitializedID, UninitializedVersion, name, phone)
}
