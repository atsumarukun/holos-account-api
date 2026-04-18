package handler_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase/dto"
	"github.com/atsumarukun/holos-account-api/test/mock/usecase"
	"github.com/atsumarukun/holos-api-pkg/errors"
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
		requestBody      []byte
		expectCode       int
		expectResponse   []byte
		setMockSessionUC func(*usecase.MockSessionUsecase)
	}{
		{
			name:           "successfully loggedin",
			requestBody:    []byte(`{"account_name":"name","password":"password"}`),
			expectCode:     http.StatusOK,
			expectResponse: fmt.Appendf(nil, `{"token":"%s"}`, sessionDTO.Token),
			setMockSessionUC: func(sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sessionDTO, nil).
					Times(1)
			},
		},
		{
			name:             "bad request",
			requestBody:      nil,
			expectCode:       http.StatusBadRequest,
			expectResponse:   []byte(`{"error":{"code":"BAD_REQUEST","message":"bad request"}}`),
			setMockSessionUC: func(*usecase.MockSessionUsecase) {},
		},
		{
			name:           "unauthenticated",
			requestBody:    []byte(`{"account_name":"name","password":"PASSWORD"}`),
			expectCode:     http.StatusUnauthorized,
			expectResponse: []byte(`{"error":{"code":"UNAUTHENTICATED","message":"unauthenticated"}}`),
			setMockSessionUC: func(sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(entity.ErrAccountPasswordIncorrect, errors.CodeUnauthenticated, "failed to verify account password")).
					Times(1)
			},
		},
		{
			name:           "internal server error",
			requestBody:    []byte(`{"account_name":"name","password":"password"}`),
			expectCode:     http.StatusInternalServerError,
			expectResponse: []byte(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}`),
			setMockSessionUC: func(sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to find account by name")).
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
			c.Request, err = http.NewRequestWithContext(ctx, "POST", "/login", bytes.NewBuffer(tt.requestBody))
			if err != nil {
				t.Error(err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionUC := usecase.NewMockSessionUsecase(ctrl)
			tt.setMockSessionUC(sessionUC)

			hdl := handler.NewSessionHandler(sessionUC)
			hdl.Login(c)

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

func TestSession_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectResponse        []byte
		expectCode            int
		setMockSessionUC      func(*usecase.MockSessionUsecase)
	}{
		{
			name:                  "successfully loggedout",
			hasAccountIDInContext: true,
			expectResponse:        nil,
			expectCode:            http.StatusNoContent,
			setMockSessionUC: func(sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Logout(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                  "account id not set",
			hasAccountIDInContext: false,
			expectResponse:        []byte(`{"error":{"code":"UNAUTHENTICATED","message":"unauthenticated"}}`),
			expectCode:            http.StatusUnauthorized,
			setMockSessionUC:      func(*usecase.MockSessionUsecase) {},
		},
		{
			name:                  "internal server error",
			hasAccountIDInContext: true,
			expectResponse:        []byte(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}`),
			expectCode:            http.StatusInternalServerError,
			setMockSessionUC: func(sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Logout(gomock.Any(), gomock.Any()).
					Return(errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to find session by account_id")).
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
			if tt.hasAccountIDInContext {
				c.Set("accountID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionUC := usecase.NewMockSessionUsecase(ctrl)
			tt.setMockSessionUC(sessionUC)

			hdl := handler.NewSessionHandler(sessionUC)
			hdl.Logout(c)

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

func TestSession_Authorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountDTO := &dto.AccountDTO{
		ID:       uuid.New(),
		Name:     "name",
		Password: "$2a$10$o7qO5pbzyAfDkBcx7Mbw9.cNCyY9V/jTjPzdSMbbwb6IixUHg3PZK",
	}

	tests := []struct {
		name                  string
		hasAccountIDInContext bool
		expectResponse        []byte
		expectCode            int
		setMockSessionUC      func(*usecase.MockSessionUsecase)
	}{
		{
			name:                  "successfully authorized",
			hasAccountIDInContext: true,
			expectResponse:        fmt.Appendf(nil, `{"id":"%s","name":"%s"}`, accountDTO.ID, accountDTO.Name),
			expectCode:            http.StatusOK,
			setMockSessionUC: func(sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Authorize(gomock.Any(), gomock.Any()).
					Return(accountDTO, nil).
					Times(1)
			},
		},
		{
			name:                  "account id not found",
			hasAccountIDInContext: false,
			expectResponse:        []byte(`{"error":{"code":"UNAUTHENTICATED","message":"unauthenticated"}}`),
			expectCode:            http.StatusUnauthorized,
			setMockSessionUC:      func(*usecase.MockSessionUsecase) {},
		},
		{
			name:                  "internal server error",
			hasAccountIDInContext: true,
			expectResponse:        []byte(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}}`),
			expectCode:            http.StatusInternalServerError,
			setMockSessionUC: func(sessionUC *usecase.MockSessionUsecase) {
				sessionUC.
					EXPECT().
					Authorize(gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(sql.ErrConnDone, errors.CodeInternalServerError, "failed to find session by token and not expired")).
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
			c.Request, err = http.NewRequestWithContext(ctx, "GET", "/authorization", http.NoBody)
			if err != nil {
				t.Error(err)
			}
			if tt.hasAccountIDInContext {
				c.Set("accountID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionUC := usecase.NewMockSessionUsecase(ctrl)
			tt.setMockSessionUC(sessionUC)

			hdl := handler.NewSessionHandler(sessionUC)
			hdl.Authorize(c)

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
