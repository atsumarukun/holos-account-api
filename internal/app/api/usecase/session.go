//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/mapper"
)

type SessionUsecase interface {
	Login(context.Context, string, string) (*dto.SessionDTO, error)
}

type sessionUsecase struct {
	transactionObj transaction.TransactionObject
	sessionRepo    repository.SessionRepository
	accountRepo    repository.AccountRepository
}

func NewSessionUsecase(
	transactionObj transaction.TransactionObject,
	sessionRepo repository.SessionRepository,
	accountRepo repository.AccountRepository,
) SessionUsecase {
	return &sessionUsecase{
		transactionObj: transactionObj,
		sessionRepo:    sessionRepo,
		accountRepo:    accountRepo,
	}
}

func (u *sessionUsecase) Login(ctx context.Context, accountName, password string) (*dto.SessionDTO, error) {
	var session *entity.Session

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		account, err := u.accountRepo.FindOneByName(ctx, accountName)
		if err != nil {
			return err
		}
		if account == nil {
			return status.ErrUnauthorized
		}

		if err := account.ComparePassword(password); err != nil {
			return err
		}

		session, err = entity.NewSession(account)
		if err != nil {
			return err
		}

		return u.sessionRepo.Save(ctx, session)
	}); err != nil {
		return nil, err
	}

	return mapper.ToSessionDTO(session), nil
}
