package model

import "github.com/google/uuid"

type AccountModel struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Password string    `db:"password"`
}
