package database

import (
	"context"
	"database/sql"
	stderr "errors"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/model"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/transformer"
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
	const errMessage = "failed to create account"

	if account == nil {
		return errors.Wrap(repository.ErrNilAccount, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToAccountModel(account)

	if _, err := driver.ExecContext(ctx, `INSERT INTO accounts (id, name, password) VALUES (?, ?, ?);`, model.ID, model.Name, model.Password); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *accountRepository) Update(ctx context.Context, account *entity.Account) error {
	const errMessage = "failed to update account"

	if account == nil {
		return errors.Wrap(repository.ErrNilAccount, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToAccountModel(account)

	if _, err := driver.ExecContext(ctx, `UPDATE accounts SET name = ?, password = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;`, model.Name, model.Password, model.ID); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *accountRepository) Delete(ctx context.Context, account *entity.Account) error {
	const errMessage = "failed to delete account"

	if account == nil {
		return errors.Wrap(repository.ErrNilAccount, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToAccountModel(account)

	if _, err := driver.ExecContext(ctx, `UPDATE accounts SET deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;`, model.ID); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *accountRepository) FindOneByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	const errMessage = "faild to find account by id"

	return r.findOne(
		ctx,
		`SELECT id, name, password FROM accounts WHERE id = ? AND deleted_at IS NULL LIMIT 1;`,
		[]any{id},
		errMessage,
	)
}

func (r *accountRepository) FindOneByName(ctx context.Context, name string) (*entity.Account, error) {
	const errMessage = "faild to find account by name"

	return r.findOne(
		ctx,
		`SELECT id, name, password FROM accounts WHERE name = ? AND deleted_at IS NULL LIMIT 1;`,
		[]any{name},
		errMessage,
	)
}

func (r *accountRepository) FindOneByNameIncludingDeleted(ctx context.Context, name string) (*entity.Account, error) {
	const errMessage = "faild to find account by name including deleted"

	return r.findOne(
		ctx,
		`SELECT id, name, password FROM accounts WHERE name = ? LIMIT 1;`,
		[]any{name},
		errMessage,
	)
}

// nolint:dupl // 集約単位のrepository実装. 集約境界を保つためrepository間での共通化は行わず重複を許容.
func (r *accountRepository) findOne(ctx context.Context, query string, args []any, errMessage string) (*entity.Account, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.AccountModel

	if err := driver.QueryRowxContext(ctx, query, args...).StructScan(&model); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return transformer.ToAccountEntity(&model), nil
}
