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
		setMockTransactionObj func(*transaction.MockTransactionObject)
		setMockAccountRepo    func(*repository.MockAccountRepository)
		setMockAccountServ    func(*service.MockAccountService)
	}{
		{
			name:                 "successfully created",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         accountDTO,
			expectError:          nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockAccountServ: func(accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
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
			setMockTransactionObj: func(*transaction.MockTransactionObject) {},
			setMockAccountRepo:    func(*repository.MockAccountRepository) {},
			setMockAccountServ:    func(*service.MockAccountService) {},
		},
		{
			name:                  "invalid password",
			inputName:             "name",
			inputPassword:         "",
			inputConfirmPassword:  "",
			expectResult:          nil,
			expectError:           status.ErrBadRequest,
			setMockTransactionObj: func(*transaction.MockTransactionObject) {},
			setMockAccountRepo:    func(*repository.MockAccountRepository) {},
			setMockAccountServ:    func(*service.MockAccountService) {},
		},
		{
			name:                 "account already exists",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          status.ErrConflict,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(*repository.MockAccountRepository) {},
			setMockAccountServ: func(accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
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
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockAccountServ: func(accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
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
			tt.setMockTransactionObj(transactionObj)

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			accountServ := service.NewMockAccountService(ctrl)
			tt.setMockAccountServ(accountServ)

			uc := usecase.NewAccountUsecase(transactionObj, accountRepo, accountServ)
			result, err := uc.Create(ctx, tt.inputName, tt.inputPassword, tt.inputConfirmPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.AccountDTO{}, "ID", "Password"),
			}
			if diff := cmp.Diff(tt.expectResult, result, opts...); diff != "" {
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
		Password: account.Password,
	}

	tests := []struct {
		name                  string
		inputID               uuid.UUID
		inputPassword         string
		inputName             string
		expectResult          *dto.AccountDTO
		expectError           error
		setMockTransactionObj func(*transaction.MockTransactionObject)
		setMockAccountRepo    func(*repository.MockAccountRepository)
		setMockAccountServ    func(*service.MockAccountService)
	}{
		{
			name:          "successfully updated",
			inputID:       account.ID,
			inputPassword: "password",
			inputName:     "update",
			expectResult:  accountDTO,
			expectError:   nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockAccountServ: func(accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "authentication failed",
			inputID:       account.ID,
			inputPassword: "PASSWORD",
			inputName:     "update",
			expectResult:  nil,
			expectError:   status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
			setMockAccountServ: func(*service.MockAccountService) {},
		},
		{
			name:          "name not changed",
			inputID:       account.ID,
			inputPassword: "password",
			inputName:     "update",
			expectResult:  accountDTO,
			expectError:   nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
			setMockAccountServ: func(*service.MockAccountService) {},
		},
		{
			name:          "invalid name",
			inputID:       account.ID,
			inputPassword: "password",
			inputName:     "",
			expectResult:  nil,
			expectError:   status.ErrBadRequest,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
			setMockAccountServ: func(*service.MockAccountService) {},
		},
		{
			name:          "find error",
			inputID:       account.ID,
			inputPassword: "password",
			inputName:     "name",
			expectResult:  nil,
			expectError:   sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAccountServ: func(*service.MockAccountService) {},
		},
		{
			name:          "account already exists",
			inputID:       account.ID,
			inputPassword: "password",
			inputName:     "name",
			expectResult:  nil,
			expectError:   status.ErrConflict,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
			setMockAccountServ: func(accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
					Return(status.ErrConflict).
					Times(1)
			},
		},
		{
			name:          "update error",
			inputID:       account.ID,
			inputPassword: "password",
			inputName:     "update",
			expectResult:  nil,
			expectError:   sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockAccountServ: func(accountServ *service.MockAccountService) {
				accountServ.
					EXPECT().
					Exists(gomock.Any(), gomock.Any()).
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
			tt.setMockTransactionObj(transactionObj)

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			accountServ := service.NewMockAccountService(ctrl)
			tt.setMockAccountServ(accountServ)

			uc := usecase.NewAccountUsecase(transactionObj, accountRepo, accountServ)
			result, err := uc.UpdateName(ctx, tt.inputID, tt.inputPassword, tt.inputName)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAccount_UpdatePassword(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}
	accountDTO := &dto.AccountDTO{
		ID:       account.ID,
		Name:     account.Name,
		Password: "$2a$10$aAjIc6dW5T07F3WzoWGnq.qGO2rMwoVAjDVeH6/t86AsIs/uIgMAG	",
	}

	tests := []struct {
		name                  string
		inputID               uuid.UUID
		inputPassword         string
		inputNewPassword      string
		inputConfirmPassword  string
		expectResult          *dto.AccountDTO
		expectError           error
		setMockTransactionObj func(*transaction.MockTransactionObject)
		setMockAccountRepo    func(*repository.MockAccountRepository)
	}{
		{
			name:                 "successfully updated",
			inputID:              account.ID,
			inputPassword:        "password",
			inputNewPassword:     "password",
			inputConfirmPassword: "password",
			expectResult:         accountDTO,
			expectError:          nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                 "authentication failed",
			inputID:              account.ID,
			inputPassword:        "PASSWORD",
			inputNewPassword:     "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
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
			name:                 "invalid password",
			inputID:              account.ID,
			inputPassword:        "password",
			inputNewPassword:     "",
			inputConfirmPassword: "",
			expectResult:         nil,
			expectError:          status.ErrBadRequest,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
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
			name:                 "find error",
			inputID:              account.ID,
			inputPassword:        "password",
			inputNewPassword:     "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
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
		{
			name:                 "update error",
			inputID:              account.ID,
			inputPassword:        "password",
			inputNewPassword:     "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
					EXPECT().
					Update(gomock.Any(), gomock.Any()).
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

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			uc := usecase.NewAccountUsecase(transactionObj, accountRepo, nil)
			result, err := uc.UpdatePassword(ctx, tt.inputID, tt.inputPassword, tt.inputNewPassword, tt.inputConfirmPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.AccountDTO{}, "Password"),
			}
			if diff := cmp.Diff(tt.expectResult, result, opts...); diff != "" {
				t.Error(diff)
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
		name                  string
		inputID               uuid.UUID
		inputPassword         string
		expectError           error
		setMockTransactionObj func(*transaction.MockTransactionObject)
		setMockAccountRepo    func(*repository.MockAccountRepository)
	}{
		{
			name:          "successfully deleted",
			inputID:       account.ID,
			inputPassword: "password",
			expectError:   nil,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
					EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "authentication failed",
			inputID:       account.ID,
			inputPassword: "PASSWORD",
			expectError:   status.ErrUnauthorized,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
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
			name:          "find error",
			inputID:       account.ID,
			inputPassword: "password",
			expectError:   sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
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
		{
			name:          "delete error",
			inputID:       account.ID,
			inputPassword: "password",
			expectError:   sql.ErrConnDone,
			setMockTransactionObj: func(transactionObj *transaction.MockTransactionObject) {
				transactionObj.
					EXPECT().
					Transaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByID(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
				accountRepo.
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

			accountRepo := repository.NewMockAccountRepository(ctrl)
			tt.setMockAccountRepo(accountRepo)

			uc := usecase.NewAccountUsecase(transactionObj, accountRepo, nil)
			if err := uc.Delete(ctx, tt.inputID, tt.inputPassword); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
