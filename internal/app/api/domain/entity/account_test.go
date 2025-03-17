package entity_test

import (
	"strings"
	"testing"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/google/uuid"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name                 string
		inputName            string
		inputPassword        string
		inputConfirmPassword string
		expectError          error
	}{
		{name: "success", inputName: "name", inputPassword: "password", inputConfirmPassword: "password", expectError: nil},
		{name: "invalid name", inputName: "", inputPassword: "password", inputConfirmPassword: "password", expectError: status.ErrBadRequest},
		{name: "invalid password", inputName: "name", inputPassword: "", inputConfirmPassword: "", expectError: status.ErrBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := entity.NewAccount(tt.inputName, tt.inputPassword, tt.inputConfirmPassword)
			if !status.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
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
		{name: "lower case only", inputName: "name", expectError: nil},
		{name: "upper case only", inputName: "NAME", expectError: nil},
		{name: "number only", inputName: "1234", expectError: nil},
		{name: "mixed lower case and upper case and number", inputName: "accountName1234", expectError: nil},
		{name: "include underscore", inputName: "account_name", expectError: nil},
		{name: "hiragana", inputName: "なまえ", expectError: status.ErrBadRequest},
		{name: "katakana", inputName: "ナマエ", expectError: status.ErrBadRequest},
		{name: "kanji", inputName: "ナマエ", expectError: status.ErrBadRequest},
		{name: "include hyphen", inputName: "account-name", expectError: status.ErrBadRequest},
		{name: "2 characters", inputName: strings.Repeat("a", 2), expectError: status.ErrBadRequest},
		{name: "3 characters", inputName: strings.Repeat("a", 3), expectError: nil},
		{name: "24 characters", inputName: strings.Repeat("a", 24), expectError: nil},
		{name: "25 characters", inputName: strings.Repeat("a", 25), expectError: status.ErrBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := account.SetName(tt.inputName); !status.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
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
		name                 string
		inputPassword        string
		inputConfirmPassword string
		expectError          error
	}{
		{name: "lower case only", inputPassword: "password", inputConfirmPassword: "password", expectError: nil},
		{name: "upper case only", inputPassword: "PASSWORD", inputConfirmPassword: "PASSWORD", expectError: nil},
		{name: "number only", inputPassword: "12345678", inputConfirmPassword: "12345678", expectError: nil},
		{name: "mixed lower case and upper case and number", inputPassword: "accountPassword1234", inputConfirmPassword: "accountPassword1234", expectError: nil},
		{name: "all symbols", inputPassword: "!@#$%^&*()_-+=[]{};:'\",.<>?/|~", inputConfirmPassword: "!@#$%^&*()_-+=[]{};:'\",.<>?/|~", expectError: nil},
		{name: "hiragana", inputPassword: "ぱすわーどぱすわーど", inputConfirmPassword: "ぱすわーどぱすわーど", expectError: status.ErrBadRequest},
		{name: "katakana", inputPassword: "パスワードパスワード", inputConfirmPassword: "パスワードパスワード", expectError: status.ErrBadRequest},
		{name: "kanji", inputPassword: "暗証番号暗証番号", inputConfirmPassword: "暗証番号暗証番号", expectError: status.ErrBadRequest},
		{name: "7 characters", inputPassword: strings.Repeat("a", 7), inputConfirmPassword: strings.Repeat("a", 7), expectError: status.ErrBadRequest},
		{name: "8 characters", inputPassword: strings.Repeat("a", 8), inputConfirmPassword: strings.Repeat("a", 8), expectError: nil},
		{name: "72 characters", inputPassword: strings.Repeat("a", 72), inputConfirmPassword: strings.Repeat("a", 72), expectError: nil},
		{name: "73 characters", inputPassword: strings.Repeat("a", 73), inputConfirmPassword: strings.Repeat("a", 73), expectError: status.ErrBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := account.SetPassword(tt.inputPassword, tt.inputConfirmPassword); !status.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if account.Password == tt.inputPassword {
					t.Error("password is not hashed")
				}
			}
		})
	}
}
