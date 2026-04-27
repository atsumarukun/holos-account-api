package builder

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
)

func ToSessionResponse(session *dto.SessionDTO) *schema.SessionResponse {
	if session == nil {
		return nil
	}

	return &schema.SessionResponse{
		Token: session.Token,
	}
}

func ToVerifiedSessionResponse(account *dto.AccountDTO) *schema.VerifiedSessionResponse {
	if account == nil {
		return nil
	}

	return &schema.VerifiedSessionResponse{
		ID:   account.ID,
		Name: account.Name,
	}
}
