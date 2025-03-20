package dto

import "github.com/google/uuid"

type AccountDTO struct {
	ID       uuid.UUID
	Name     string
	Password string
}
