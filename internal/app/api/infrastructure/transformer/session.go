package transformer

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/model"
)

func ToSessionModel(session *entity.Session) *model.SessionModel {
	return &model.SessionModel{
		AccountID: session.AccountID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	}
}

func ToSessionEntity(session *model.SessionModel) *entity.Session {
	return &entity.Session{
		AccountID: session.AccountID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	}
}
