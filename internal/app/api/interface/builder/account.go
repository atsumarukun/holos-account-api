package builder

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
)

func ToAccountResponse(account *dto.AccountDTO) *schema.AccountResponse {
	if account == nil {
		return nil
	}

	return &schema.AccountResponse{
		Name: account.Name,
	}
}
