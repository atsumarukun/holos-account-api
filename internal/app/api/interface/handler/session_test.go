package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/test/mock/usecase"
)

func TestSession_Logtin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sessionDTO := &dto.SessionDTO{
		AccountID: uuid.New(),
		Token:     "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	tests := []struct {
		name             string
		requestJSON      []byte
		expectCode       int
		expectResponse   map[string]any
		setMockSessionUC func(context.Context, *usecase.MockSessionUsecase)
	}{
		{
			name:           "success",
			requestJSON:    []byte(`{"account_name": "name", "password": "password"}`),
			expectCode:     http.StatusOK,
			expectResponse: map[string]any{"token": "1Ty1HKTPKTt8xEi-_3HTbWf2SCHOdqOS"},
			setMockSessionUC: func(ctx context.Context, sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Login(ctx, gomock.Any(), gomock.Any()).
					Return(sessionDTO, nil).
					Times(1)
			},
		},
		{
			name:             "invalid request",
			requestJSON:      nil,
			expectCode:       http.StatusBadRequest,
			expectResponse:   map[string]any{"message": "bad request"},
			setMockSessionUC: func(context.Context, *usecase.MockSessionUsecase) {},
		},
		{
			name:           "login error",
			requestJSON:    []byte(`{"account_name": "name", "password": "password"}`),
			expectCode:     http.StatusUnauthorized,
			expectResponse: map[string]any{"message": "unauthorized"},
			setMockSessionUC: func(ctx context.Context, sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Login(ctx, gomock.Any(), gomock.Any()).
					Return(nil, status.ErrUnauthorized).
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/login", bytes.NewBuffer(tt.requestJSON))
			if err != nil {
				t.Error(err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionUC := usecase.NewMockSessionUsecase(ctrl)
			tt.setMockSessionUC(ctx, sessionUC)

			hdl := handler.NewSessionHandler(sessionUC)
			hdl.Login(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			var response map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(response, tt.expectResponse); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestSession_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		isSetAccountID   bool
		expectCode       int
		setMockSessionUC func(context.Context, *usecase.MockSessionUsecase)
	}{
		{
			name:           "success",
			isSetAccountID: true,
			expectCode:     http.StatusNoContent,
			setMockSessionUC: func(ctx context.Context, sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Logout(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:             "account id not found",
			isSetAccountID:   false,
			expectCode:       http.StatusInternalServerError,
			setMockSessionUC: func(context.Context, *usecase.MockSessionUsecase) {},
		},
		{
			name:           "logout faild",
			isSetAccountID: true,
			expectCode:     http.StatusUnauthorized,
			setMockSessionUC: func(ctx context.Context, sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Logout(ctx, gomock.Any()).
					Return(status.ErrUnauthorized).
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
			c.Request, err = http.NewRequestWithContext(ctx, "DELETE", "/logout", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			if tt.isSetAccountID {
				c.Set("accountID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionUC := usecase.NewMockSessionUsecase(ctrl)
			tt.setMockSessionUC(ctx, sessionUC)

			hdl := handler.NewSessionHandler(sessionUC)
			hdl.Logout(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}
		})
	}
}
