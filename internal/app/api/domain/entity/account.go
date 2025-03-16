package entity

import (
	"errors"

	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID
	Name     string
	Password string
}

func NewAccount(name string, password string, confirmPassword string) (*Account, error) {
	return nil, errors.New("not implemented")
}

func RestoreAccount(id uuid.UUID, name string, password string) *Account {
	return &Account{
		ID:       id,
		Name:     name,
		Password: password,
	}
}

func (a *Account) SetName(name string) error {
	return errors.New("not implemented")
}

func (a *Account) SetPassword(password string, confirmPassword string) error {
	return errors.New("not implemented")
}
