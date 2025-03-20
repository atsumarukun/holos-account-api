package database_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	mockDatabase "github.com/atsumarukun/holos-account-api/test/mock/database"
)

func TestSession_Create(t *testing.T) {
	session := &entity.Session{
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
				mock.ExpectExec(regexp.QuoteMeta(`REPLACE sessions (account_id, token, expires_at) VALUES (?, ?, ?);`)).
					WithArgs(session.AccountID, session.Token, session.ExpiresAt).
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
			name:         "replace error",
			inputSession: session,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`REPLACE sessions (account_id, token, expires_at) VALUES (?, ?, ?);`)).
					WithArgs(session.AccountID, session.Token, session.ExpiresAt).
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
			if err := repo.Save(t.Context(), tt.inputSession); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSession_Delete(t *testing.T) {
	session := &entity.Session{
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
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sessions WHERE account_id = ?;`)).
					WithArgs(session.AccountID).
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
			name:         "delete error",
			inputSession: session,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sessions WHERE account_id = ?;`)).
					WithArgs(session.AccountID).
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
			if err := repo.Delete(t.Context(), tt.inputSession); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
