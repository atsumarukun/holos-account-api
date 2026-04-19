//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	stderr "errors"

	"github.com/atsumarukun/holos-api-pkg/errors"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
)

var ErrAccountNameAlreadyInUse = stderr.New("account name already in use")

type AccountService interface {
	Exists(context.Context, *entity.Account) error
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) Exists(ctx context.Context, account *entity.Account) error {
	acc, err := s.accountRepo.FindOneByNameIncludingDeleted(ctx, account.Name)
	if err != nil {
		return err
	}
	if acc != nil {
		return errors.Wrap(ErrAccountNameAlreadyInUse, errors.CodeDuplicate, "account already exists")
	}
	return nil
}
