package database_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	mockDatabase "github.com/atsumarukun/holos-account-api/test/mock/database"
)

func TestAccount_Create(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name         string
		inputAccount *entity.Account
		expectError  error
		setMockDB    func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputAccount: account,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts (id, name, password) VALUES (?, ?, ?);`)).
					WithArgs(account.ID, account.Name, account.Password).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:         "account is nil",
			inputAccount: nil,
			expectError:  status.ErrInternal,
			setMockDB:    func(sqlmock.Sqlmock) {},
		},
		{
			name:         "insert error",
			inputAccount: account,
			expectError:  status.FromError(sql.ErrConnDone),
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts (id, name, password) VALUES (?, ?, ?);`)).
					WithArgs(account.ID, account.Name, account.Password).
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

			r := database.NewDBAccountRepository(db)
			if err := r.Create(t.Context(), tt.inputAccount); !status.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAccount_Update(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name         string
		inputAccount *entity.Account
		expectError  error
		setMockDB    func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputAccount: account,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE accounts SET name = ?, password = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(account.Name, account.Password, account.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:         "account is nil",
			inputAccount: nil,
			expectError:  status.ErrInternal,
			setMockDB:    func(sqlmock.Sqlmock) {},
		},
		{
			name:         "update error",
			inputAccount: account,
			expectError:  status.FromError(sql.ErrConnDone),
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE accounts SET name = ?, password = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(account.Name, account.Password, account.ID).
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

			r := database.NewDBAccountRepository(db)
			if err := r.Update(t.Context(), tt.inputAccount); !status.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAccount_Delete(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name         string
		inputAccount *entity.Account
		expectError  error
		setMockDB    func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputAccount: account,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE accounts SET deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(account.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:         "account is nil",
			inputAccount: nil,
			expectError:  status.ErrInternal,
			setMockDB:    func(sqlmock.Sqlmock) {},
		},
		{
			name:         "delete error",
			inputAccount: account,
			expectError:  status.FromError(sql.ErrConnDone),
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE accounts SET deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(account.ID).
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

			r := database.NewDBAccountRepository(db)
			if err := r.Delete(t.Context(), tt.inputAccount); !status.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAccount_FindOneByID(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name         string
		inputID      uuid.UUID
		expectResult *entity.Account
		expectError  error
		setMockDB    func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputID:      account.ID,
			expectResult: account,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(account.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"}).AddRow(account.ID, account.Name, account.Password)).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputID:      account.ID,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(account.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"})).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputID:      account.ID,
			expectResult: nil,
			expectError:  status.FromError(sql.ErrConnDone),
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(account.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)
			defer db.Close()

			tt.setMockDB(mock)

			r := database.NewDBAccountRepository(db)
			result, err := r.FindOneByID(t.Context(), tt.inputID)
			if !status.Is(err, tt.expectError) {
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

func TestAccount_FindOneByNameIncludingDeleted(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name         string
		inputName    string
		expectResult *entity.Account
		expectError  error
		setMockDB    func(mock sqlmock.Sqlmock)
	}{
		{
			name:         "success",
			inputName:    "name",
			expectResult: account,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE name = ? LIMIT 1;`)).
					WithArgs("name").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"}).AddRow(account.ID, account.Name, account.Password)).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputName:    "name",
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE name = ? LIMIT 1;`)).
					WithArgs("name").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"})).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputName:    "name",
			expectResult: nil,
			expectError:  status.FromError(sql.ErrConnDone),
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE name = ? LIMIT 1;`)).
					WithArgs("name").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"})).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := mockDatabase.NewMockDatabase(t)
			defer db.Close()

			tt.setMockDB(mock)

			r := database.NewDBAccountRepository(db)
			result, err := r.FindOneByNameIncludingDeleted(t.Context(), tt.inputName)
			if !status.Is(err, tt.expectError) {
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
