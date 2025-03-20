package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Token     string
	ExpiresAt time.Time
}

func NewSession(account *Account) (*Session, error) {
	return nil, errors.New("not implemented")
}

func RestoreSession(id uuid.UUID, accountID uuid.UUID, token string, expiresAt time.Time) *Session {
	return nil
}

func (s *Session) GenerateToken() error {
	return errors.New("not implemented")
}
