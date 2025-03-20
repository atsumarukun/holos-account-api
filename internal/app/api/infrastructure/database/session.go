package database

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewDBSessionRepository(db *sqlx.DB) repository.SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Create(ctx context.Context, session *entity.Session) error {
	return errors.New("not implemented")
}

func (r *sessionRepository) Update(ctx context.Context, session *entity.Session) error {
	return errors.New("not implemented")
}

func (r *sessionRepository) FindOneByAccountID(ctx context.Context, accountID uuid.UUID) (*entity.Session, error) {
	return nil, errors.New("not implemented")
}
