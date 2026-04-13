//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	stderr "errors"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
)

var ErrNilAccount = stderr.New("account must not be nil")

type AccountRepository interface {
	Create(context.Context, *entity.Account) error
	Update(context.Context, *entity.Account) error
	Delete(context.Context, *entity.Account) error
	FindOneByID(context.Context, uuid.UUID) (*entity.Account, error)
	FindOneByName(context.Context, string) (*entity.Account, error)
	FindOneByNameIncludingDeleted(context.Context, string) (*entity.Account, error)
}
