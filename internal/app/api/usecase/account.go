//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/mapper"
)

type AccountUsecase interface {
	Create(context.Context, string, string, string) (*dto.AccountDTO, error)
	UpdateName(context.Context, uuid.UUID, string, string) (*dto.AccountDTO, error)
	UpdatePassword(context.Context, uuid.UUID, string, string) (*dto.AccountDTO, error)
	Delete(context.Context, uuid.UUID) error
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

func (u *accountUsecase) UpdateName(ctx context.Context, id uuid.UUID, password, name string) (*dto.AccountDTO, error) {
	var account *entity.Account

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		var err error
		account, err = u.accountRepo.FindOneByID(ctx, id)
		if err != nil {
			return err
		}
		if account == nil {
			return status.ErrUnauthorized
		}

		if err := account.ComparePassword(password); err != nil {
			return err
		}

		if account.Name == name {
			return nil
		}

		if err := account.SetName(name); err != nil {
			return err
		}

		if err := u.accountServ.Exists(ctx, account); err != nil {
			return err
		}

		return u.accountRepo.Update(ctx, account)
	}); err != nil {
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}

func (u *accountUsecase) UpdatePassword(ctx context.Context, id uuid.UUID, password, confirmPassword string) (*dto.AccountDTO, error) {
	var account *entity.Account

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		var err error
		account, err = u.accountRepo.FindOneByID(ctx, id)
		if err != nil {
			return err
		}
		if account == nil {
			return status.ErrUnauthorized
		}

		if err := account.SetPassword(password, confirmPassword); err != nil {
			return err
		}

		return u.accountRepo.Update(ctx, account)
	}); err != nil {
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}

func (u *accountUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		account, err := u.accountRepo.FindOneByID(ctx, id)
		if err != nil {
			return err
		}
		if account == nil {
			return status.ErrUnauthorized
		}

		return u.accountRepo.Delete(ctx, account)
	})
}
