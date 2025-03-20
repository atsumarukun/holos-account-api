package entity

import (
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
)

type Account struct {
	ID       uuid.UUID
	Name     string
	Password string
}

func NewAccount(name, password, confirmPassword string) (*Account, error) {
	var account Account

	if err := account.generateID(); err != nil {
		return nil, err
	}
	if err := account.SetName(name); err != nil {
		return nil, err
	}
	if err := account.SetPassword(password, confirmPassword); err != nil {
		return nil, err
	}

	return &account, nil
}

func RestoreAccount(id uuid.UUID, name, password string) *Account {
	return &Account{
		ID:       id,
		Name:     name,
		Password: password,
	}
}

func (a *Account) SetName(name string) error {
	if len(name) < 3 || 24 < len(name) {
		return status.ErrBadRequest
	}
	if matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name); err != nil {
		return err
	} else if !matched {
		return status.ErrBadRequest
	}
	a.Name = name
	return nil
}

func (a *Account) SetPassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return status.ErrBadRequest
	}
	if len(password) < 8 || 72 < len(password) {
		return status.ErrBadRequest
	}
	if matched, err := regexp.MatchString(`^[A-Za-z0-9!@#$%^&*()_\-+=\[\]{};:'",.<>?/\\|~]*$`, password); err != nil {
		return err
	} else if !matched {
		return status.ErrBadRequest
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashed)
	return nil
}

func (a *Account) generateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}
