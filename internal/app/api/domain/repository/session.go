//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	stderr "errors"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
)

var ErrNilSession = stderr.New("session must not be nil")

type SessionRepository interface {
	Save(context.Context, *entity.Session) error
	Delete(context.Context, *entity.Session) error
	FindOneByAccountID(context.Context, uuid.UUID) (*entity.Session, error)
	FindOneByTokenAndNotExpired(context.Context, string) (*entity.Session, error)
}
