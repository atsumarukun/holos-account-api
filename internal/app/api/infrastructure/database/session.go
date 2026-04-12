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

type sessionRepository struct {
	db *sqlx.DB
}

func NewDBSessionRepository(db *sqlx.DB) repository.SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Save(ctx context.Context, session *entity.Session) error {
	const errMessage = "failed to save session"

	if session == nil {
		return errors.Wrap(repository.ErrRequiredSession, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToSessionModel(session)

	if _, err := driver.ExecContext(ctx, `REPLACE sessions (account_id, token, expires_at) VALUES (?, ?, ?);`, model.AccountID, model.Token, model.ExpiresAt); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, session *entity.Session) error {
	const errMessage = "failed to delete session"

	if session == nil {
		return errors.Wrap(repository.ErrRequiredSession, errors.CodeInternalServerError, errMessage)
	}

	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToSessionModel(session)

	if _, err := driver.ExecContext(ctx, `DELETE FROM sessions WHERE account_id = ?;`, model.AccountID); err != nil {
		return errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return nil
}

func (r *sessionRepository) FindOneByAccountID(ctx context.Context, accountID uuid.UUID) (*entity.Session, error) {
	const errMessage = "faild to find session by account_id"

	driver := transaction.GetDriver(ctx, r.db)
	var model model.SessionModel

	if err := driver.QueryRowxContext(ctx, `SELECT account_id, token, expires_at FROM sessions WHERE account_id = ?;`, accountID).StructScan(&model); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return transformer.ToSessionEntity(&model), nil
}

func (r *sessionRepository) FindOneByTokenAndNotExpired(ctx context.Context, token string) (*entity.Session, error) {
	const errMessage = "faild to find session by tolen and not expired"

	driver := transaction.GetDriver(ctx, r.db)
	var model model.SessionModel

	if err := driver.QueryRowxContext(ctx, `SELECT account_id, token, expires_at FROM sessions WHERE token = ? AND expires_at > NOW(6);`, token).StructScan(&model); err != nil {
		if stderr.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, errors.CodeInternalServerError, errMessage)
	}

	return transformer.ToSessionEntity(&model), nil
}
