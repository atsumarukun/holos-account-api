package database_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	mockDatabase "github.com/atsumarukun/holos-account-api/test/mock/database"
)

func TestSession_Create(t *testing.T) {
	session := &entity.Session{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	tests := []struct {
		name         string
		inputSession *entity.Session
		expectError  error
		setMockDB    func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputSession: session,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sessions (id, account_id, token, expires_at) VALUES (?, ?, ?, ?);`)).
					WithArgs(session.ID, session.AccountID, session.Token, session.ExpiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:         "session is nil",
			inputSession: nil,
			expectError:  status.ErrInternal,
			setMockDB:    func(mock sqlmock.Sqlmock) {},
		},

		{
			name:         "insert error",
			inputSession: session,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sessions (id, account_id, token, expires_at) VALUES (?, ?, ?, ?);`)).
					WithArgs(session.ID, session.AccountID, session.Token, session.ExpiresAt).
					WillReturnResult(nil).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)
			defer db.Close()

			tt.setMockDB(mock)

			repo := database.NewDBSessionRepository(db)
			if err := repo.Create(t.Context(), tt.inputSession); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSession_Update(t *testing.T) {
	session := &entity.Session{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	tests := []struct {
		name         string
		inputSession *entity.Session
		expectError  error
		setMockDB    func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputSession: session,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE sessions SET token = ?, expires_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(session.Token, session.ExpiresAt, session.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:         "session is nil",
			inputSession: nil,
			expectError:  status.ErrInternal,
			setMockDB:    func(mock sqlmock.Sqlmock) {},
		},

		{
			name:         "update error",
			inputSession: session,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE sessions SET token = ?, expires_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(session.Token, session.ExpiresAt, session.ID).
					WillReturnResult(nil).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)
			defer db.Close()

			tt.setMockDB(mock)

			repo := database.NewDBSessionRepository(db)
			if err := repo.Update(t.Context(), tt.inputSession); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSession_FindOneByAccountID(t *testing.T) {
	session := &entity.Session{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	tests := []struct {
		name           string
		inputAccountID uuid.UUID
		expectResult   *entity.Session
		expectError    error
		setMockDB      func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			expectResult: session,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, token, expires_at FROM sessions WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(session.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "token", "expires_at"}).AddRow(session.ID, session.AccountID, session.Token, session.ExpiresAt)).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			expectResult: nil,
			expectError:  status.ErrInternal,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, token, expires_at FROM sessions WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(session.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "token", "expires_at"})).
					WillReturnError(nil)
			},
		},

		{
			name:         "find error",
			expectResult: session,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, account_id, token, expires_at FROM sessions WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(session.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "token", "expires_at"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)
			defer db.Close()

			tt.setMockDB(mock)

			repo := database.NewDBSessionRepository(db)
			result, err := repo.FindOneByAccountID(t.Context(), tt.inputAccountID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
