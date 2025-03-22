package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler"
	"github.com/atsumarukun/holos-account-api/internal/app/api/pkg/status"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/test/mock/usecase"
)

func TestAccount_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountDTO := &dto.AccountDTO{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name             string
		requestJSON      []byte
		expectCode       int
		expectResponse   map[string]any
		setMockAccountUC func(context.Context, *usecase.MockAccountUsecase)
	}{
		{
			name:           "success",
			requestJSON:    []byte(`{"name": "name", "password": "password", "confirm_password": "password"}`),
			expectCode:     http.StatusCreated,
			expectResponse: map[string]any{"name": "name"},
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:             "invalid request",
			requestJSON:      nil,
			expectCode:       http.StatusBadRequest,
			expectResponse:   map[string]any{"message": "bad request"},
			setMockAccountUC: func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:           "create error",
			requestJSON:    []byte(`{"name": "name", "password": "password", "confirm_password": "password"}`),
			expectCode:     http.StatusConflict,
			expectResponse: map[string]any{"message": "conflict"},
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, status.ErrConflict).
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/accounts", bytes.NewBuffer(tt.requestJSON))
			if err != nil {
				t.Error(err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accountUC := usecase.NewMockAccountUsecase(ctrl)
			tt.setMockAccountUC(ctx, accountUC)

			hdl := handler.NewAccountHandler(accountUC)
			hdl.Create(c)

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

func TestAccount_UpdateName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountDTO := &dto.AccountDTO{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name             string
		requestJSON      []byte
		isSetAccountID   bool
		expectCode       int
		expectResponse   map[string]any
		setMockAccountUC func(context.Context, *usecase.MockAccountUsecase)
	}{
		{
			name:           "success",
			requestJSON:    []byte(`{"name": "name"}`),
			isSetAccountID: true,
			expectCode:     http.StatusOK,
			expectResponse: map[string]any{"name": "name"},
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdateName(ctx, gomock.Any(), gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:             "invalid request",
			requestJSON:      nil,
			isSetAccountID:   true,
			expectCode:       http.StatusBadRequest,
			expectResponse:   map[string]any{"message": "bad request"},
			setMockAccountUC: func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:             "account id not found",
			requestJSON:      []byte(`{"name": "name"}`),
			isSetAccountID:   false,
			expectCode:       http.StatusInternalServerError,
			expectResponse:   map[string]any{"message": "internal server error"},
			setMockAccountUC: func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:           "update error",
			requestJSON:    []byte(`{"name": "name"}`),
			isSetAccountID: true,
			expectCode:     http.StatusConflict,
			expectResponse: map[string]any{"message": "conflict"},
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdateName(ctx, gomock.Any(), gomock.Any()).
					Return(nil, status.ErrConflict).
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
			c.Request, err = http.NewRequestWithContext(ctx, "PUT", "/accounts/name", bytes.NewBuffer(tt.requestJSON))
			if err != nil {
				t.Error(err)
			}
			if tt.isSetAccountID {
				c.Set("accountID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accountUC := usecase.NewMockAccountUsecase(ctrl)
			tt.setMockAccountUC(ctx, accountUC)

			hdl := handler.NewAccountHandler(accountUC)
			hdl.UpdateName(c)

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

func TestAccount_UpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountDTO := &dto.AccountDTO{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name             string
		requestJSON      []byte
		isSetAccountID   bool
		expectCode       int
		expectResponse   map[string]any
		setMockAccountUC func(context.Context, *usecase.MockAccountUsecase)
	}{
		{
			name:           "success",
			requestJSON:    []byte(`{"password": "password", "confirm_password": "password"}`),
			isSetAccountID: true,
			expectCode:     http.StatusOK,
			expectResponse: map[string]any{"name": "name"},
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdatePassword(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:             "invalid request",
			requestJSON:      nil,
			isSetAccountID:   true,
			expectCode:       http.StatusBadRequest,
			expectResponse:   map[string]any{"message": "bad request"},
			setMockAccountUC: func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:             "account id not found",
			requestJSON:      []byte(`{"password": "password", "confirm_password": "password"}`),
			isSetAccountID:   false,
			expectCode:       http.StatusInternalServerError,
			expectResponse:   map[string]any{"message": "internal server error"},
			setMockAccountUC: func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:           "update error",
			requestJSON:    []byte(`{"password": "password", "confirm_password": "password"}`),
			isSetAccountID: true,
			expectCode:     http.StatusUnauthorized,
			expectResponse: map[string]any{"message": "unauthorized"},
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdatePassword(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "PUT", "/accounts/password", bytes.NewBuffer(tt.requestJSON))
			if err != nil {
				t.Error(err)
			}
			if tt.isSetAccountID {
				c.Set("accountID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accountUC := usecase.NewMockAccountUsecase(ctrl)
			tt.setMockAccountUC(ctx, accountUC)

			hdl := handler.NewAccountHandler(accountUC)
			hdl.UpdatePassword(c)

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
