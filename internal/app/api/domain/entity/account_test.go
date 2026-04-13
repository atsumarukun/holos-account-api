package entity_test

import (
	stderr "errors"
	"strings"
	"testing"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name              string
		inputName         string
		inputPassword     string
		inputConfirmation string
		expectError       error
	}{
		{name: "successfully initialized", inputName: "name", inputPassword: "password", inputConfirmation: "password", expectError: nil},
		{name: "invalid name", inputName: "", inputPassword: "password", inputConfirmation: "password", expectError: entity.ErrAccountNameInvalidLength},
		{name: "invalid password", inputName: "name", inputPassword: "", inputConfirmation: "", expectError: entity.ErrAccountPasswordInvalidLength},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := entity.NewAccount(tt.inputName, tt.inputPassword, tt.inputConfirmation)
			if !stderr.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err != nil {
				if _, ok := err.(interface {
					Code() errors.ErrorCode
					Message() string
				}); !ok {
					t.Errorf("error is not wrapped")
				}
			}

			if tt.expectError == nil {
				if account == nil {
					t.Error("account is nil")
				} else {
					if account.ID == uuid.Nil {
						t.Error("id is not set")
					}
					if account.Name == "" {
						t.Error("name is not set")
					}
					if account.Password == "" {
						t.Error("password is not set")
					}
					if len(account.Password) != 60 {
						t.Error("password length is not 60")
					}
				}
			}
		})
	}
}

func TestAccount_SetName(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name        string
		inputName   string
		expectError error
	}{
		{name: "mixed lower case and upper case and number", inputName: "accountName1234", expectError: nil},
		{name: "include underscore", inputName: "account_name", expectError: nil},
		{name: "include hyphen", inputName: "account-name", expectError: entity.ErrAccountNameInvalidChars},
		{name: "full-width characters", inputName: "アカウント名", expectError: entity.ErrAccountNameInvalidChars},
		{name: "2 characters", inputName: strings.Repeat("a", 2), expectError: entity.ErrAccountNameInvalidLength},
		{name: "3 characters", inputName: strings.Repeat("a", 3), expectError: nil},
		{name: "24 characters", inputName: strings.Repeat("a", 24), expectError: nil},
		{name: "25 characters", inputName: strings.Repeat("a", 25), expectError: entity.ErrAccountNameInvalidLength},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := account.SetName(tt.inputName)
			if !stderr.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err != nil {
				if _, ok := err.(interface {
					Code() errors.ErrorCode
					Message() string
				}); !ok {
					t.Errorf("error is not wrapped")
				}
			}
		})
	}
}

func TestAccount_SetPassword(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name              string
		inputPassword     string
		inputConfirmation string
		expectError       error
	}{
		{name: "lower case only", inputPassword: "password", inputConfirmation: "password", expectError: nil},
		{name: "upper case only", inputPassword: "PASSWORD", inputConfirmation: "PASSWORD", expectError: nil},
		{name: "number only", inputPassword: "12345678", inputConfirmation: "12345678", expectError: nil},
		{name: "mixed lower case and upper case and number", inputPassword: "accountPassword1234", inputConfirmation: "accountPassword1234", expectError: nil},
		{name: "all symbols", inputPassword: "!@#$%^&*()_-+=[]{};:'\",.<>?/|~", inputConfirmation: "!@#$%^&*()_-+=[]{};:'\",.<>?/|~", expectError: nil},
		{name: "full-width characters", inputPassword: "認証パスワードぱすわーど", inputConfirmation: "認証パスワードぱすわーど", expectError: entity.ErrAccountPasswordInvalidChars},
		{name: "7 characters", inputPassword: strings.Repeat("a", 7), inputConfirmation: strings.Repeat("a", 7), expectError: entity.ErrAccountPasswordInvalidLength},
		{name: "8 characters", inputPassword: strings.Repeat("a", 8), inputConfirmation: strings.Repeat("a", 8), expectError: nil},
		{name: "72 characters", inputPassword: strings.Repeat("a", 72), inputConfirmation: strings.Repeat("a", 72), expectError: nil},
		{name: "73 characters", inputPassword: strings.Repeat("a", 73), inputConfirmation: strings.Repeat("a", 73), expectError: entity.ErrAccountPasswordInvalidLength},
		{name: "dose not matched", inputPassword: "password", inputConfirmation: "PASSWORD", expectError: entity.ErrAccountPasswordMismatch},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := account.SetPassword(tt.inputPassword, tt.inputConfirmation)
			if !stderr.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err != nil {
				if _, ok := err.(interface {
					Code() errors.ErrorCode
					Message() string
				}); !ok {
					t.Errorf("error is not wrapped")
				}
			}

			if tt.expectError == nil {
				if account.Password == tt.inputPassword {
					t.Error("password is not hashed")
				}
			}
		})
	}
}

func TestAccount_VerifyPassword(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name          string
		inputPassword string
		expectError   error
	}{
		{
			name:          "successfully verified",
			inputPassword: "password",
			expectError:   nil,
		},
		{
			name:          "faild",
			inputPassword: "PASSWORD",
			expectError:   entity.ErrAccountPasswordIncorrect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := account.VerifyPassword(tt.inputPassword)
			if !stderr.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err != nil {
				if _, ok := err.(interface {
					Code() errors.ErrorCode
					Message() string
				}); !ok {
					t.Errorf("error is not wrapped")
				}
			}
		})
	}
}
