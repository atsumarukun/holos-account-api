package entity

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Token     string
	ExpiresAt time.Time
}

func NewSession(account *Account) (*Session, error) {
	var session Session

	if err := session.generateID(); err != nil {
		return nil, err
	}
	if err := session.setAccount(account); err != nil {
		return nil, err
	}
	if err := session.GenerateToken(); err != nil {
		return nil, err
	}

	return &session, nil
}

func RestoreSession(id uuid.UUID, accountID uuid.UUID, token string, expiresAt time.Time) *Session {
	return &Session{
		ID:        id,
		AccountID: accountID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

func (s *Session) GenerateToken() error {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf)
	if 32 < len(token) {
		return status.ErrInternal
	}
	s.Token = token
	s.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	return nil
}

func (s *Session) generateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (s *Session) setAccount(account *Account) error {
	if account == nil {
		return status.ErrInternal
	}
	s.AccountID = account.ID
	return nil
}
