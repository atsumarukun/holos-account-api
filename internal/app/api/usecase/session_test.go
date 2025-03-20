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
		setMockTransactionObj func(context.Context, *transaction.MockTransactionObject)
		setMockSessionRepo    func(context.Context, *repository.MockSessionRepository)
		setMockAccountRepo    func(context.Context, *repository.MockAccountRepository)
	}{
		{
			name:             "success",
			inputAccountName: "name",
			inputPassword:    "password",
			expectResult:     sessionDTO,
			expectError:      nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(ctx context.Context, sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(ctx, gomock.Any()).
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
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(context.Context, *repository.MockSessionRepository) {},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(ctx, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:             "password dose not matched",
			inputAccountName: "name",
			inputPassword:    "PASSWORD",
			expectResult:     nil,
			expectError:      status.ErrUnauthorized,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(context.Context, *repository.MockSessionRepository) {},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(ctx, gomock.Any()).
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
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(context.Context, *repository.MockSessionRepository) {},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(ctx, gomock.Any()).
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
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(ctx context.Context, sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					Save(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByName(ctx, gomock.Any()).
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
			tt.setMockTransactionObj(ctx, transactionObj)

			sessionRepo := repository.NewMockSessionRepository(ctrl)
			tt.setMockSessionRepo(ctx, sessionRepo)

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(ctx, accountRepo)

			uc := usecase.NewSessionUsecase(transactionObj, sessionRepo, accountRepo)
			result, err := uc.Login(ctx, tt.inputAccountName, tt.inputPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.SessionDTO{}, "Token", "ExpiresAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
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
		setMockTransactionObj func(context.Context, *transaction.MockTransactionObject)
		setMockSessionRepo    func(context.Context, *repository.MockSessionRepository)
	}{
		{
			name:           "success",
			inputAccountID: session.AccountID,
			expectError:    nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(ctx context.Context, sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(ctx, gomock.Any()).
					Return(session, nil).
					Times(1)
				sessionRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "session not found",
			inputAccountID: session.AccountID,
			expectError:    status.ErrUnauthorized,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(ctx context.Context, sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(ctx, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:           "find session error",
			inputAccountID: session.AccountID,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(ctx context.Context, sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(ctx, gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:           "delete session error",
			inputAccountID: session.AccountID,
			expectError:    sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockSessionRepo: func(ctx context.Context, sessionRepo *repository.MockSessionRepository) {
				sessionRepo.
					EXPECT().
					FindOneByAccountID(ctx, gomock.Any()).
					Return(session, nil).
					Times(1)
				sessionRepo.
					EXPECT().
					Delete(ctx, gomock.Any()).
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
			tt.setMockTransactionObj(ctx, transactionObj)

			sessionRepo := repository.NewMockSessionRepository(ctrl)
			tt.setMockSessionRepo(ctx, sessionRepo)

			uc := usecase.NewSessionUsecase(transactionObj, sessionRepo, nil)
			if err := uc.Logout(ctx, tt.inputAccountID); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
