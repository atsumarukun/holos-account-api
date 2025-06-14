package database_test

import (
	"database/sql"
	"errors"
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
			name:         "successfully inserted",
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
			expectError:  sql.ErrConnDone,
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

			repo := database.NewDBAccountRepository(db)
			if err := repo.Create(t.Context(), tt.inputAccount); !errors.Is(err, tt.expectError) {
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
			name:         "successfully updated",
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
			expectError:  sql.ErrConnDone,
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

			repo := database.NewDBAccountRepository(db)
			if err := repo.Update(t.Context(), tt.inputAccount); !errors.Is(err, tt.expectError) {
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
			name:         "successfully deleted",
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
			expectError:  sql.ErrConnDone,
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

			repo := database.NewDBAccountRepository(db)
			if err := repo.Delete(t.Context(), tt.inputAccount); !errors.Is(err, tt.expectError) {
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
			name:         "successfully found",
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
			expectError:  sql.ErrConnDone,
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

			repo := database.NewDBAccountRepository(db)
			result, err := repo.FindOneByID(t.Context(), tt.inputID)
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

func TestAccount_FindOneByName(t *testing.T) {
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
			name:         "successfully found",
			inputName:    "name",
			expectResult: account,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE name = ? AND deleted_at IS NULL LIMIT 1;`)).
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE name = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs("name").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"})).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputName:    "name",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, password FROM accounts WHERE name = ? AND deleted_at IS NULL LIMIT 1;`)).
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

			repo := database.NewDBAccountRepository(db)
			result, err := repo.FindOneByName(t.Context(), tt.inputName)
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
			name:         "successfully found",
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
			expectError:  sql.ErrConnDone,
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

			repo := database.NewDBAccountRepository(db)
			result, err := repo.FindOneByNameIncludingDeleted(t.Context(), tt.inputName)
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
