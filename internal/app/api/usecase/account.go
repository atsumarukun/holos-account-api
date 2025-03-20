//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/mapper"
)

type AccountUsecase interface {
	Create(context.Context, string, string, string) (*dto.AccountDTO, error)
}

type accountUsecase struct {
	transactionObj transaction.TransactionObject
	accountRepo    repository.AccountRepository
	accountServ    service.AccountService
}

func NewAccountUsecase(
	transactionObj transaction.TransactionObject,
	accountRepo repository.AccountRepository,
	accountServ service.AccountService,
) AccountUsecase {
	return &accountUsecase{
		transactionObj: transactionObj,
		accountRepo:    accountRepo,
		accountServ:    accountServ,
	}
}

func (u *accountUsecase) Create(ctx context.Context, name, password, confirmPassword string) (*dto.AccountDTO, error) {
	account, err := entity.NewAccount(name, password, confirmPassword)
	if err != nil {
		return nil, err
	}

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		if err := u.accountServ.Exists(ctx, account); err != nil {
			return err
		}

		return u.accountRepo.Create(ctx, account)
	}); err != nil {
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}
