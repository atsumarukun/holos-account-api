package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/middleware"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/test/mock/usecase"
)

func TestAuthentication_Authenticate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountDTO := &dto.AccountDTO{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name                string
		authorizationHeader string
		expectResult        uuid.UUID
		setMockSessionUC    func(context.Context, *usecase.MockSessionUsecase)
	}{
		{
			name:                "success",
			authorizationHeader: "Session 1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult:        accountDTO.ID,
			setMockSessionUC: func(ctx context.Context, sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Authenticate(ctx, gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:                "invalid authorization header",
			authorizationHeader: "",
			expectResult:        uuid.Nil,
			setMockSessionUC:    func(context.Context, *usecase.MockSessionUsecase) {},
		},
		{
			name:                "account not found",
			authorizationHeader: "Session 1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
			expectResult:        uuid.Nil,
			setMockSessionUC: func(ctx context.Context, sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Authenticate(ctx, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			var err error
			c.Request, err = http.NewRequestWithContext(ctx, "DELETE", "/logout", nil)
			if err != nil {
				t.Error(err)
			}
			c.Request.Header.Add("Authorization", tt.authorizationHeader)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionUC := usecase.NewMockSessionUsecase(ctrl)
			tt.setMockSessionUC(ctx, sessionUC)

			mw := middleware.NewAuthenticationMiddleware(sessionUC)
			mw.Authenticate(c)

			accountID, exists := c.Get("accountID")
			if exists && tt.expectResult == uuid.Nil {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, accountID)
			} else {
				result, _ := accountID.(uuid.UUID)
				if diff := cmp.Diff(result, tt.expectResult); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
