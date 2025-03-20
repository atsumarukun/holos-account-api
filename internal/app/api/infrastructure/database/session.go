package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/model"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/transformer"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
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
	if session == nil {
		return status.ErrInternal
	}
	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToSessionModel(session)
	if _, err := driver.ExecContext(ctx, `INSERT INTO sessions (id, account_id, token, expires_at) VALUES (?, ?, ?, ?);`, model.ID, model.AccountID, model.Token, model.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) Update(ctx context.Context, session *entity.Session) error {
	if session == nil {
		return status.ErrInternal
	}
	driver := transaction.GetDriver(ctx, r.db)
	model := transformer.ToSessionModel(session)
	if _, err := driver.ExecContext(ctx, `UPDATE sessions SET token = ?, expires_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;`, model.Token, model.ExpiresAt, model.ID); err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) FindOneByAccountID(ctx context.Context, accountID uuid.UUID) (*entity.Session, error) {
	driver := transaction.GetDriver(ctx, r.db)
	var model model.SessionModel
	if err := driver.QueryRowxContext(ctx, `SELECT id, account_id, token, expires_at FROM sessions WHERE account_id = ? AND deleted_at IS NULL LIMIT 1;`, accountID).StructScan(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transformer.ToSessionEntity(&model), nil
}
