package entity

import (
	"crypto/rand"
	"encoding/base64"
	stderr "errors"
	"time"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-api-pkg/errors"
)

var (
	ErrSessionTokenInvalidLength = stderr.New("token must be 32 characters")
	ErrSessionNilAccount         = stderr.New("account must not be nil")
)

type Session struct {
	AccountID uuid.UUID
	Token     string
	ExpiresAt time.Time
}

func NewSession(account *Account) (*Session, error) {
	var session Session

	if err := session.setAccount(account); err != nil {
		return nil, err
	}
	if err := session.GenerateToken(); err != nil {
		return nil, err
	}

	return &session, nil
}

func RestoreSession(accountID uuid.UUID, token string, expiresAt time.Time) *Session {
	return &Session{
		AccountID: accountID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

func (s *Session) GenerateToken() error {
	const errMessage = "failed to generate token"

	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf)
	if 32 != len(token) {
		return errors.Wrap(ErrSessionTokenInvalidLength, errors.CodeInternalServerError, errMessage)
	}

	s.Token = token
	s.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)

	return nil
}

func (s *Session) setAccount(account *Account) error {
	if account == nil {
		return errors.Wrap(ErrSessionNilAccount, errors.CodeInternalServerError, "failed to set session account")
	}
	s.AccountID = account.ID
	return nil
}
