package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/test/mock/domain/repository"
	"github.com/atsumarukun/holos-account-api/test/mock/domain/repository/pkg/transaction"
)

func TestSession_Login(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}
	sessionDTO := &dto.SessionDTO{
		AccountID: account.ID,
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	tests := []struct {
		name                  string
		inputAccountName      string
		inputPassword         string
		expectResult          *dto.SessionDTO
		expectError           error
		setMockTransactionObj func(*transaction.MockTransactionObject)
		setMockSessionRepo    func(*repository.MockSessionRepository)
		setMockAccountRepo    func(*repository.MockAccountRepository)
	}{
		{
			name:             "successfully loggedin",
			inputAccountName: "name",
			inputPassword:    "password",
			expectResult:     sessionDTO,
			expectError:      nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
		},
		{
			name:             "account not found",
			inputAccountName: "name",
			inputPassword:    "password",
			expectResult:     nil,
			expectError:      status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(*repository.MockSessionRepository) {},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:             "authentication failed",
			inputAccountName: "name",
			inputPassword:    "PASSWORD",
			expectResult:     nil,
			expectError:      status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(*repository.MockSessionRepository) {},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
		},
		{
			name:             "find account error",
			inputAccountName: "name",
			inputPassword:    "password",
			expectResult:     nil,
			expectError:      sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(*repository.MockSessionRepository) {},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:             "save session error",
			inputAccountName: "name",
			inputPassword:    "password",
			expectResult:     nil,
			expectError:      sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			transactionObj := transaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(transactionObj)

			sessionRepo := repository.NewMockSessionRepository(ctrl)
			tt.setMockSessionRepo(sessionRepo)

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			uc := usecase.NewSessionUsecase(transactionObj, sessionRepo, accountRepo)
			result, err := uc.Login(ctx, tt.inputAccountName, tt.inputPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.SessionDTO{}, "Token", "ExpiresAt"),
			}
			if diff := cmp.Diff(tt.expectResult, result, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestSession_Logout(t *testing.T) {
	session := &entity.Session{
		AccountID: uuid.New(),
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	tests := []struct {
		name                  string
		inputAccountID        uuid.UUID
		expectError           error
		setMockTransactionObj func(*transaction.MockTransactionObject)
		setMockSessionRepo    func(*repository.MockSessionRepository)
	}{
		{
			name:           "successfully loggedout",
			inputAccountID: session.AccountID,
			expectError:    nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(gomock.Any(), gomock.Any()).
					Return(session, nil).
					Times(1)
				sessionRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "session not found",
			inputAccountID: session.AccountID,
			expectError:    status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:           "find session error",
			inputAccountID: session.AccountID,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:           "delete session error",
			inputAccountID: session.AccountID,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(gomock.Any(), gomock.Any()).
					Return(session, nil).
					Times(1)
				sessionRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			transactionObj := transaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(transactionObj)

			sessionRepo := repository.NewMockSessionRepository(ctrl)
			tt.setMockSessionRepo(sessionRepo)

			uc := usecase.NewSessionUsecase(transactionObj, sessionRepo, nil)
			if err := uc.Logout(ctx, tt.inputAccountID); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestSession_Authenticate(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}
	session := &entity.Session{
		AccountID: account.ID,
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	accountDTO := &dto.AccountDTO{
		ID:       account.ID,
		Name:     account.Name,
		Password: account.Password,
	}

	tests := []struct {
		name                  string
		inputToken            string
		expectResult          *dto.AccountDTO
		expectError           error
		setMockTransactionObj func(*transaction.MockTransactionObject)
		setMockSessionRepo    func(*repository.MockSessionRepository)
		setMockAccountRepo    func(*repository.MockAccountRepository)
	}{
		{
			name:         "successfully authenticated",
			inputToken:   "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult: accountDTO,
			expectError:  nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByTokenAndNotExpired(gomock.Any(), gomock.Any()).
					Return(session, nil).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
		},
		{
			name:         "session not found",
			inputToken:   "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult: nil,
			expectError:  status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByTokenAndNotExpired(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			setMockAccountRepo: func(*repository.MockAccountRepository) {},
		},
		{
			name:         "account not found",
			inputToken:   "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult: nil,
			expectError:  status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByTokenAndNotExpired(gomock.Any(), gomock.Any()).
					Return(session, nil).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find session error",
			inputToken:   "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByTokenAndNotExpired(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAccountRepo: func(*repository.MockAccountRepository) {},
		},
		{
			name:         "find account error",
			inputToken:   "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByTokenAndNotExpired(gomock.Any(), gomock.Any()).
					Return(session, nil).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			transactionObj := transaction.NewMockTransactionObject(ctrl)
			tt.setMockTransactionObj(transactionObj)

			sessionRepo := repository.NewMockSessionRepository(ctrl)
			tt.setMockSessionRepo(sessionRepo)

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			uc := usecase.NewSessionUsecase(transactionObj, sessionRepo, accountRepo)
			result, err := uc.Authenticate(ctx, tt.inputToken)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestSession_Authorize(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}
	accountDTO := &dto.AccountDTO{
		ID:       account.ID,
		Name:     account.Name,
		Password: account.Password,
	}

	tests := []struct {
		name               string
		inputAccountID     uuid.UUID
		expectResult       *dto.AccountDTO
		expectError        error
		setMockAccountRepo func(*repository.MockAccountRepository)
	}{
		{
			name:           "successfully authorized",
			inputAccountID: account.ID,
			expectResult:   accountDTO,
			expectError:    nil,
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
		},
		{
			name:           "account not found",
			inputAccountID: account.ID,
			expectResult:   nil,
			expectError:    status.ErrUnauthorized,
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:           "find account error",
			inputAccountID: account.ID,
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := t.Context()

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			uc := usecase.NewSessionUsecase(nil, nil, accountRepo)
			result, err := uc.Authorize(ctx, tt.inputAccountID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
