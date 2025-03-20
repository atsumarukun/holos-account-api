//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
)

type SessionRepository interface {
	Create(context.Context, *entity.Session) error
	Update(context.Context, *entity.Session) error
	FindOneByAccountID(context.Context, uuid.UUID) (*entity.Session, error)
}
