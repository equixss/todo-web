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
	nameLen := len([]rune(u.Name))
	if nameLen < 3 || nameLen > 100 {
		return fmt.Errorf("invalid name length: %d: %w", nameLen, core_errors.ErrInvalidArgument)
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

type UserPatch struct {
	Name  Nullable[string]
	Phone Nullable[string]
}

func NewUserPatch(
	name Nullable[string],
	phone Nullable[string],
) UserPatch {
	return UserPatch{
		Name:  name,
		Phone: phone,
	}
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}

	tmp := *u
	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}
	if patch.Phone.Set {
		tmp.Phone = patch.Phone.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("invalid patch user: %w", err)
	}

	*u = tmp
	return nil
}

func (p *UserPatch) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf("'Name' can't be patched to null:%w", core_errors.ErrInvalidArgument)
	}
	return nil
}
