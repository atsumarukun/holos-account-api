package service_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/test/mock/domain/repository"
)

func TestAccount_Exists(t *testing.T) {
	account := &entity.Account{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name               string
		inputAccount       *entity.Account
		expectError        error
		setMockAccountRepo func(*repository.MockAccountRepository)
	}{
		{
			name:         "not exists",
			inputAccount: account,
			expectError:  nil,
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByNameIncludingDeleted(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "exists",
			inputAccount: account,
			expectError:  status.ErrConflict,
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByNameIncludingDeleted(gomock.Any(), gomock.Any()).
					Return(account, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputAccount: account,
			expectError:  sql.ErrConnDone,
			setMockAccountRepo: func(accountRepo *repository.MockAccountRepository) {
				accountRepo.
					EXPECT().
					FindOneByNameIncludingDeleted(gomock.Any(), gomock.Any()).
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

			serv := service.NewAccountService(accountRepo)
			if err := serv.Exists(ctx, tt.inputAccount); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
