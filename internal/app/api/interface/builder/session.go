package builder

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/schema"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
)

func ToSessionResponse(session *dto.SessionDTO) *schema.SessionResponse {
	return &schema.SessionResponse{
		Token: session.Token,
	}
}
