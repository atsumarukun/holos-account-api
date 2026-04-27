//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	stderr "errors"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/mapper"
)

var ErrSessionNotFound = stderr.New("session not found")

type SessionUsecase interface {
	Create(context.Context, string, string) (*dto.SessionDTO, error)
	Delete(context.Context, uuid.UUID) error
	Verify(context.Context, string) (*dto.AccountDTO, error)
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

func (u *sessionUsecase) Create(ctx context.Context, accountName, password string) (*dto.SessionDTO, error) {
	var session *entity.Session

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		account, err := u.accountRepo.FindOneByName(ctx, accountName)
		if err != nil {
			return err
		}
		if account == nil {
			return errors.Wrap(ErrAccountNotFound, errors.CodeUnauthenticated, "failed to create session")
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

func (u *sessionUsecase) Delete(ctx context.Context, accountID uuid.UUID) error {
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

func (u *sessionUsecase) Verify(ctx context.Context, token string) (*dto.AccountDTO, error) {
	var account *entity.Account

	if err := u.transactionObj.Transaction(ctx, func(ctx context.Context) error {
		session, err := u.sessionRepo.FindOneByTokenAndNotExpired(ctx, token)
		if err != nil {
			return err
		}
		if session == nil {
			return errors.Wrap(ErrSessionNotFound, errors.CodeUnauthenticated, "failed to verify")
		}

		account, err = u.accountRepo.FindOneByID(ctx, session.AccountID)
		if err != nil {
			return err
		}
		if account == nil {
			return errors.Wrap(ErrAccountNotFound, errors.CodeUnauthenticated, "failed to verify")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.ToAccountDTO(account), nil
}
