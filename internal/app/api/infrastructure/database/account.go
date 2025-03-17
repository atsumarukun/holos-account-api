package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
)

type accountRepository struct {
	db *sqlx.DB
}

func NewDBAccountRepository(db *sqlx.DB) repository.AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Create(ctx context.Context, account *entity.Account) error {
	return errors.New("not implemented")
}

func (r *accountRepository) Update(ctx context.Context, account *entity.Account) error {
	return errors.New("not implemented")
}

func (r *accountRepository) Delete(ctx context.Context, account *entity.Account) error {
	return errors.New("not implemented")
}

func (r *accountRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	return nil, errors.New("not implemented")
}

func (r *accountRepository) FindOneByNameIncludingDeleted(ctx context.Context, name string) (*entity.Account, error) {
	return nil, errors.New("not implemented")
}
