package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atsumarukun/holos-api-pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler"
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
			requestBody:    []byte(`{"name":"name","password":"password","confirm_password":"password"}`),
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
			name:             "bad request",
			requestBody:      nil,
			expectCode:       http.StatusBadRequest,
			expectResponse:   []byte(`{"error":{"code":"BAD_REQUEST","message":"bad request"}}`),
			setMockAccountUC: func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:           "duplicate",
			requestBody:    []byte(`{"name":"name","password":"password","confirm_password":"password"}`),
			expectCode:     http.StatusConflict,
			expectResponse: []byte(`{"error":{"code":"DUPLICATE","message":"account name already in use"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(service.ErrAccountNameAlreadyInUse, errors.CodeDuplicate, "account already exists")).
					Times(1)
			},
		},
		{
			name:           "invalid input",
			requestBody:    []byte(`{"name":"名前","password":"password","confirm_password":"password"}`),
			expectCode:     http.StatusUnprocessableEntity,
			expectResponse: []byte(`{"error":{"code":"INVALID_INPUT","message":"account name contains invalid characters"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(entity.ErrAccountNameInvalidChars, errors.CodeInvalidInput, "failed to set account name")).
					Times(1)
			},
		},
		{
			name:           "internal server error",
			requestBody:    []byte(`{"name":"name","password":"password","confirm_password":"password"}`),
			expectCode:     http.StatusInternalServerError,
			expectResponse: []byte(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to create account")).
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
			requestBody:           []byte(`{"password":"password","name": "name"}`),
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
			name:                  "bad request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"error":{"code":"BAD_REQUEST","message":"bad request"}}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "account id not set",
			requestBody:           []byte(`{"name": "name"}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusUnauthorized,
			expectResponse:        []byte(`{"error":{"code":"UNAUTHENTICATED","message":"unauthenticated"}}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "duplicate",
			requestBody:           []byte(`{"password":"password","name":"name"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusConflict,
			expectResponse:        []byte(`{"error":{"code":"DUPLICATE","message":"account name already in use"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdateName(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(service.ErrAccountNameAlreadyInUse, errors.CodeDuplicate, "account already exists")).
					Times(1)
			},
		},
		{
			name:                  "invalid input",
			requestBody:           []byte(`{"password":"password","name":"名前"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusUnprocessableEntity,
			expectResponse:        []byte(`{"error":{"code":"INVALID_INPUT","message":"account name contains invalid characters"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdateName(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(entity.ErrAccountNameInvalidChars, errors.CodeInvalidInput, "failed to set account name")).
					Times(1)
			},
		},
		{
			name:                  "internal server error",
			requestBody:           []byte(`{"password":"password","name":"name"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdateName(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to update account")).
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
			requestBody:           []byte(`{"password":"password","new_password":"password","confirm_password":"password"}`),
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
			name:                  "bad request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"error":{"code":"BAD_REQUEST","message":"bad request"}}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "account id not set",
			requestBody:           []byte(`{"password":"password","new_password":"password","confirm_password":"password"}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusUnauthorized,
			expectResponse:        []byte(`{"error":{"code":"UNAUTHENTICATED","message":"unauthenticated"}}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "invalid input",
			requestBody:           []byte(`{"password":"password","new_password":"パスワード","confirm_password":"パスワード"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusUnprocessableEntity,
			expectResponse:        []byte(`{"error":{"code":"INVALID_INPUT","message":"password contains invalid characters"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdatePassword(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(entity.ErrAccountPasswordInvalidChars, errors.CodeInvalidInput, "failed to set account password")).
					Times(1)
			},
		},
		{
			name:                  "internal server error",
			requestBody:           []byte(`{"password":"password","new_password":"password","confirm_password":"password"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					UpdatePassword(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to update account")).
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
			requestBody:           []byte(`{"password":"password"}`),
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
			name:                  "bad request",
			requestBody:           nil,
			hasAccountIDInContext: true,
			expectCode:            http.StatusBadRequest,
			expectResponse:        []byte(`{"error":{"code":"BAD_REQUEST","message":"bad request"}}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "account id not found",
			requestBody:           []byte(`{"password":"password"}`),
			hasAccountIDInContext: false,
			expectCode:            http.StatusUnauthorized,
			expectResponse:        []byte(`{"error":{"code":"UNAUTHENTICATED","message":"unauthenticated"}}`),
			setMockAccountUC:      func(context.Context, *usecase.MockAccountUsecase) {},
		},
		{
			name:                  "internal server error",
			requestBody:           []byte(`{"password":"password"}`),
			hasAccountIDInContext: true,
			expectCode:            http.StatusInternalServerError,
			expectResponse:        []byte(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}`),
			setMockAccountUC: func(ctx context.Context, accountUC *usecase.MockAccountUsecase) {
				accountUC.
					EXPECT().
					Delete(ctx, gomock.Any(), gomock.Any()).
					Return(errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to find account by id")).
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
