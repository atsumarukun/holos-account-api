//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	stderr "errors"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/mapper"
	"github.com/atsumarukun/holos-api-pkg/errors"
)

var ErrSessionNotFound = stderr.New("session not found")

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
			return errors.Wrap(ErrAccountNotFound, errors.CodeUnauthenticated, "failed to login")
		}

		if err := account.VerifyPassword(password); err != nil {
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
			return nil
		}

		return u.sessionRepo.Delete(ctx, session)
	})
}

func (u *sessionUsecase) Authenticate(ctx context.Context, token string) (*dto.AccountDTO, error) {
	var account *entity.Account

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		session, err := u.sessionRepo.FindOneByTokenAndNotExpired(ctx, token)
		if err != nil {
			return err
		}
		if session == nil {
			return errors.Wrap(ErrSessionNotFound, errors.CodeUnauthenticated, "failed to authenticate")
		}

		account, err = u.accountRepo.FindOneByID(ctx, session.AccountID)
		if err != nil {
			return err
		}
		if account == nil {
			return errors.Wrap(ErrAccountNotFound, errors.CodeUnauthenticated, "failed to authenticate")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}

func (u *sessionUsecase) Authorize(ctx context.Context, accountID uuid.UUID) (*dto.AccountDTO, error) {
	account, err := u.accountRepo.FindOneByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.Wrap(ErrAccountNotFound, errors.CodeUnauthenticated, "failed to authorize")
	}

	return mapper.ToAccountDTO(account), nil
}
