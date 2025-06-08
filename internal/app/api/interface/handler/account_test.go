package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
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
		requestBody      []byte
		expectCode       int
		expectResponse   []byte
		setMockAccountUC func(context.Context, *usecase.MockAccountUsecase)
	}{
		{
			name:           "successfully created",
			requestBody:    []byte(`{"name": "name", "password": "password", "confirm_password": "password"}`),
			expectCode:     http.StatusCreated,
			expectResponse: fmt.Appendf(nil, `{"name":"%s"}`, accountDTO.Name),
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
			requestBody:      nil,
			expectCode:       http.StatusBadRequest,
			expectResponse:   []byte(`{"message":"bad request"}`),
			setMockAccountUC: func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:           "create error",
			requestBody:    []byte(`{"name": "name", "password": "password", "confirm_password": "password"}`),
			expectCode:     http.StatusConflict,
			expectResponse: []byte(`{"message":"conflict"}`),
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/accounts", bytes.NewBuffer(tt.requestBody))
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

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
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
		name                  string
		requestBody           []byte
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockAccountUC      func(context.Context, *usecase.MockAccountUsecase)
	}{
		{
			name:                  "successfully updated",
			requestBody:           []byte(`{"password": "password", "name": "name"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"name":"%s"}`, accountDTO.Name),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdateName(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "invalid request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"message":"bad request"}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "account id not set",
			requestBody:           []byte(`{"name": "name"}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "update error",
			requestBody:           []byte(`{"password": "password", "name": "name"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusConflict,
			expectResponse:        []byte(`{"message":"conflict"}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdateName(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
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
			c.Request, err = http.NewRequestWithContext(ctx, "PUT", "/accounts/name", bytes.NewBuffer(tt.requestBody))
			if err != nil {
				t.Error(err)
			}
			if tt.hasAccountIDInContext {
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

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
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
		name                  string
		requestBody           []byte
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockAccountUC      func(context.Context, *usecase.MockAccountUsecase)
	}{
		{
			name:                  "successfully updated",
			requestBody:           []byte(`{"password": "password", "new_password": "password", "confirm_password": "password"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusOK,
			expectResponse:        fmt.Appendf(nil, `{"name":"%s"}`, accountDTO.Name),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdatePassword(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "invalid request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"message":"bad request"}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "account id not set",
			requestBody:           []byte(`{"password": "password", "confirm_password": "password"}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "update error",
			requestBody:           []byte(`{"password": "password", "new_password": "password", "confirm_password": "password"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdatePassword(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
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
			c.Request, err = http.NewRequestWithContext(ctx, "PUT", "/accounts/password", bytes.NewBuffer(tt.requestBody))
			if err != nil {
				t.Error(err)
			}
			if tt.hasAccountIDInContext {
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

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAccount_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                  string
		requestBody           []byte
		hasAccountIDInContext bool
		expectCode            int
		expectResponse        []byte
		setMockAccountUC      func(context.Context, *usecase.MockAccountUsecase)
	}{
		{
			name:                  "successfully deleted",
			requestBody:           []byte(`{"password": "password"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusNoContent,
			expectResponse:        nil,
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Delete(ctx, gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                  "invalid request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"message":"bad request"}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "account id not found",
			requestBody:           []byte(`{"password": "password"}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "delete error",
			requestBody:           []byte(`{"password": "password"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"message":"internal server error"}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Delete(ctx, gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
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
			c.Request, err = http.NewRequestWithContext(ctx, "DELETE", "/accounts", bytes.NewBuffer(tt.requestBody))
			if err != nil {
				t.Error(err)
			}
			if tt.hasAccountIDInContext {
				c.Set("accountID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accountUC := usecase.NewMockAccountUsecase(ctrl)
			tt.setMockAccountUC(ctx, accountUC)

			hdl := handler.NewAccountHandler(accountUC)
			hdl.Delete(c)

			c.Writer.WriteHeaderNow()

			if w.Code != tt.expectCode {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectCode, w.Code)
			}

			if diff := cmp.Diff(tt.expectResponse, w.Body.Bytes()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
