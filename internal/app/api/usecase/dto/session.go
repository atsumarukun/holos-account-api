package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionDTO struct {
	AccountID uuid.UUID
	Token     string
	ExpiresAt time.Time
}
