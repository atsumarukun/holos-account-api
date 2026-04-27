package mapper

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
)

func ToSessionDTO(session *entity.Session) *dto.SessionDTO {
	if session == nil {
		return nil
	}

	return &dto.SessionDTO{
		AccountID: session.AccountID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	}
}
