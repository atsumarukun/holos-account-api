package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

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
	"github.com/atsumarukun/holos-account-api/test/mock/domain/service"
)

func TestAccount_Create(t *testing.T) {
	accountDTO := &dto.AccountDTO{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name                  string
		inputName             string
		inputPassword         string
		inputConfirmPassword  string
		expectResult          *dto.AccountDTO
		expectError           error
		setMockTransactionObj func(context.Context, *transaction.MockTransactionObject)
		setMockAccountRepo    func(context.Context, *repository.MockAccountRepository)
		setMockAccountServ    func(context.Context, *service.MockAccountService)
	}{
		{
			name:                 "success",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         accountDTO,
			expectError:          nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockAccountServ: func(ctx context.Context, accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                  "invalid name",
			inputName:             "",
			inputPassword:         "password",
			inputConfirmPassword:  "password",
			expectResult:          nil,
			expectError:           status.ErrBadRequest,
			setMockTransactionObj: func(context.Context, *transaction.MockTransactionObject) {},
			setMockAccountRepo:    func(context.Context, *repository.MockAccountRepository) {},
			setMockAccountServ:    func(context.Context, *service.MockAccountService) {},
		},
		{
			name:                  "invalid password",
			inputName:             "name",
			inputPassword:         "",
			inputConfirmPassword:  "",
			expectResult:          nil,
			expectError:           status.ErrBadRequest,
			setMockTransactionObj: func(context.Context, *transaction.MockTransactionObject) {},
			setMockAccountRepo:    func(context.Context, *repository.MockAccountRepository) {},
			setMockAccountServ:    func(context.Context, *service.MockAccountService) {},
		},
		{
			name:                 "account already exists",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          status.ErrConflict,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(context.Context, *repository.MockAccountRepository) {},
			setMockAccountServ: func(ctx context.Context, accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(status.ErrConflict).
					Times(1)
			},
		},
		{
			name:                 "create error",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockAccountServ: func(ctx context.Context, accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
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

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(ctx, accountRepo)

			accountServ := service.NewMockAccountService(ctrl)
			tt.setMockAccountServ(ctx, accountServ)

			uc := usecase.NewAccountUsecase(transactionObj, accountRepo, accountServ)
			result, err := uc.Create(ctx, tt.inputName, tt.inputPassword, tt.inputConfirmPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.AccountDTO{}, "ID", "Password"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAccount_UpdateName(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}
	accountDTO := &dto.AccountDTO{
		ID:       account.ID,
		Name:     "update",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name                  string
		inputID               uuid.UUID
		inputName             string
		expectResult          *dto.AccountDTO
		expectError           error
		setMockTransactionObj func(context.Context, *transaction.MockTransactionObject)
		setMockAccountRepo    func(context.Context, *repository.MockAccountRepository)
		setMockAccountServ    func(context.Context, *service.MockAccountService)
	}{
		{
			name:         "success",
			inputID:      account.ID,
			inputName:    "update",
			expectResult: accountDTO,
			expectError:  nil,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(ctx, gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockAccountServ: func(ctx context.Context, accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "invalid name",
			inputID:      account.ID,
			inputName:    "",
			expectResult: nil,
			expectError:  status.ErrBadRequest,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(ctx, gomock.Any()).
					Return(account, nil).
					Times(1)
			},
			setMockAccountServ: func(context.Context, *service.MockAccountService) {},
		},
		{
			name:         "account already exists",
			inputID:      account.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  status.ErrConflict,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(ctx, gomock.Any()).
					Return(account, nil).
					Times(1)
			},
			setMockAccountServ: func(ctx context.Context, accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(status.ErrConflict).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputID:      account.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(ctx, gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAccountServ: func(context.Context, *service.MockAccountService) {},
		},
		{
			name:         "update error",
			inputID:      account.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObj: func(ctx context.Context, transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(ctx context.Context, accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(ctx, gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
					EXPECT().
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockAccountServ: func(ctx context.Context, accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(ctx, gomock.Any()).
					Return(nil).
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

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(ctx, accountRepo)

			accountServ := service.NewMockAccountService(ctrl)
			tt.setMockAccountServ(ctx, accountServ)

			uc := usecase.NewAccountUsecase(transactionObj, accountRepo, accountServ)
			result, err := uc.UpdateName(ctx, tt.inputID, tt.inputName)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
