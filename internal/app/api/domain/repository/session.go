//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../test/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/google/uuid"
)

type SessionRepository interface {
	Save(context.Context, *entity.Session) error
	Delete(context.Context, *entity.Session) error
	FindOneByAccountID(context.Context, uuid.UUID) (*entity.Session, error)
}
