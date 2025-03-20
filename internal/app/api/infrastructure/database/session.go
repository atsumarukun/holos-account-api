package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/model"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/transformer"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
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
	if session == nil {
		return status.ErrInternal
	}
	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToSessionModel(session)
	if _, err := driver.ExecContext(ctx, `REPLACE sessions (account_id, token, expires_at) VALUES (?, ?, ?);`, model.AccountID, model.Token, model.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, session *entity.Session) error {
	if session == nil {
		return status.ErrInternal
	}
	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToSessionModel(session)
	if _, err := driver.ExecContext(ctx, `DELETE FROM sessions WHERE account_id = ?;`, model.AccountID); err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) FindOneByAccountID(ctx context.Context, accountID uuid.UUID) (*entity.Session, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.SessionModel
	if err := driver.QueryRowxContext(ctx, `SELECT account_id, token, expires_at FROM sessions WHERE id = ?;`, accountID).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transformer.ToSessionEntity(&model), nil
}
