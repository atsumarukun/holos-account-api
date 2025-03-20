package model

import (
	"time"

	"github.com/google/uuid"
)

type SessionModel struct {
	AccountID uuid.UUID `db:"account_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
