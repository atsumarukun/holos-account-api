package entity

import (
	stderr "errors"
	"regexp"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAccountNameInvalidLength     = stderr.New("account name must be between 3 and 24 characters")
	ErrAccountNameInvalidChars      = stderr.New("account name contains invalid characters")
	ErrAccountPasswordMismatch      = stderr.New("passwords do not match")
	ErrAccountPasswordInvalidLength = stderr.New("password must be between 8 and 72 characters")
	ErrAccountPasswordInvalidChars  = stderr.New("password contains invalid characters")
	ErrAccountPasswordIncorrect     = stderr.New("password is incorrect")
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
	const errMessage = "failed to set account name"

	if len(name) < 3 || 24 < len(name) {
		return errors.Wrap(ErrAccountNameInvalidLength, errors.CodeInvalidInput, errMessage)
	}

	if matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	} else if !matched {
		return errors.Wrap(ErrAccountNameInvalidChars, errors.CodeInvalidInput, errMessage)
	}

	a.Name = name

	return nil
}

func (a *Account) SetPassword(password, confirmation string) error {
	const errMessage = "failed to set account password"

	if password != confirmation {
		return errors.Wrap(ErrAccountPasswordMismatch, errors.CodeInvalidInput, errMessage)
	}

	if len(password) < 8 || 72 < len(password) {
		return errors.Wrap(ErrAccountPasswordInvalidLength, errors.CodeInvalidInput, errMessage)
	}

	if matched, err := regexp.MatchString(`^[A-Za-z0-9!@#$%^&*()_\-+=\[\]{};:'",.<>?/\\|~]*$`, password); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	} else if !matched {
		return errors.Wrap(ErrAccountPasswordInvalidChars, errors.CodeInvalidInput, errMessage)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	a.Password = string(hashed)

	return nil
}

func (a *Account) VerifyPassword(password string) error {
	const errMessage = "failed to verify account password"

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		if stderr.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.Wrap(ErrAccountPasswordIncorrect, errors.CodeUnauthenticated, errMessage)
		}
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}
	return nil
}

func (a *Account) generateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, "failed to generate account id")
	}
	a.ID = id
	return nil
}
