//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package usecase

import (
	"context"
	"errors"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/repository/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
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
	return nil, errors.New("not implemented")
}
