//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/mapper"
)

type SessionUsecase interface {
	Login(context.Context, string, string) (*dto.SessionDTO, error)
	Logout(context.Context, uuid.UUID) error
	Authenticate(context.Context, string) (*dto.AccountDTO, error)
	Authorize(context.Context, uuid.UUID) (*dto.AccountDTO, error)
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

func (u *sessionUsecase) Logout(ctx context.Context, accountID uuid.UUID) error {
	return u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		session, err := u.sessionRepo.FindOneByAccountID(ctx, accountID)
		if err != nil {
			return err
		}
		if session == nil {
			return status.ErrUnauthorized
		}

		return u.sessionRepo.Delete(ctx, session)
	})
}

func (u *sessionUsecase) Authenticate(ctx context.Context, token string) (*dto.AccountDTO, error) {
	var account *entity.Account

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		session, err := u.sessionRepo.FindOneByToken(ctx, token)
		if err != nil {
			return err
		}
		if session == nil {
			return status.ErrUnauthorized
		}

		account, err = u.accountRepo.FindOneByID(ctx, session.AccountID)
		if err != nil {
			return err
		}
		if account == nil {
			return status.ErrUnauthorized
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}

func (u *sessionUsecase) Authorize(ctx context.Context, accountID uuid.UUID) (*dto.AccountDTO, error) {
	return nil, errors.New("not implemented")
}
