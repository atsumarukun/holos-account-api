package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/model"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/transformer"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
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
	if account == nil {
		return status.ErrInternal
	}
	driver := getDriver(ctx, r.db)
	model := transformer.ToAccountModel(account)
	if _, err := driver.ExecContext(ctx, `INSERT INTO accounts (id, name, password) VALUES (?, ?, ?);`, model.ID, model.Name, model.Password); err != nil {
		return status.FromError(err)
	}
	return nil
}

func (r *accountRepository) Update(ctx context.Context, account *entity.Account) error {
	if account == nil {
		return status.ErrInternal
	}
	driver := getDriver(ctx, r.db)
	model := transformer.ToAccountModel(account)
	if _, err := driver.ExecContext(ctx, `UPDATE accounts SET name = ?, password = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;`, model.Name, model.Password, model.ID); err != nil {
		return status.FromError(err)
	}
	return nil
}

func (r *accountRepository) Delete(ctx context.Context, account *entity.Account) error {
	if account == nil {
		return status.ErrInternal
	}
	driver := getDriver(ctx, r.db)
	model := transformer.ToAccountModel(account)
	if _, err := driver.ExecContext(ctx, `UPDATE accounts SET deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;`, model.ID); err != nil {
		return status.FromError(err)
	}
	return nil
}

func (r *accountRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	driver := getDriver(ctx, r.db)
	var model model.AccountModel
	if err := driver.QueryRowxContext(ctx, `SELECT id, name, password FROM accounts WHERE id = ? AND deleted_at IS NULL LIMIT 1;`, id).StructScan(&model); err != nil {
		if status.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, status.FromError(err)
	}
	return transformer.ToAccountEntity(&model)
}

func (r *accountRepository) FindOneByNameIncludingDeleted(ctx context.Context, name string) (*entity.Account, error) {
	driver := getDriver(ctx, r.db)
	var model model.AccountModel
	if err := driver.QueryRowxContext(ctx, `SELECT id, name, password FROM accounts WHERE name = ? LIMIT 1;`, name).StructScan(&model); err != nil {
		if status.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, status.FromError(err)
	}
	return transformer.ToAccountEntity(&model)
}
