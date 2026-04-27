package transformer

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database/model"
)

func ToSessionModel(session *entity.Session) *model.SessionModel {
	if session == nil {
		return nil
	}

	return &model.SessionModel{
		AccountID: session.AccountID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	}
}

func ToSessionEntity(session *model.SessionModel) *entity.Session {
	if session == nil {
		return nil
	}

	return entity.RestoreSession(session.AccountID, session.Token, session.ExpiresAt)
}
