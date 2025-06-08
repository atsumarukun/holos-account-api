package entity_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
)

func TestNewSession(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name         string
		inputAccount *entity.Account
		expectError  error
	}{
		{name: "successfully initialized", inputAccount: account, expectError: nil},
		{name: "account is nil", inputAccount: nil, expectError: status.ErrInternal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := entity.NewSession(tt.inputAccount)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if tt.expectError == nil {
				if session == nil {
					t.Error("session is nil")
				} else {
					if session.AccountID == uuid.Nil {
						t.Error("account_id is not set")
					}
					if session.Token == "" {
						t.Error("token is not set")
					}
					if session.ExpiresAt.Before(time.Now()) {
						t.Error("invalid expires_at")
					}
				}
			}
		})
	}
}

func TestSession_GenerateToken(t *testing.T) {
	session := &entity.Session{
		AccountID: uuid.New(),
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	tests := []struct {
		name        string
		expectError error
	}{
		{name: "successfully generated", expectError: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := session.Token

			if err := session.GenerateToken(); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if session.Token == old {
				t.Error("token has not been updated")
			}
			if len(session.Token) != 32 {
				t.Error("invalid token")
			}
		})
	}
}
