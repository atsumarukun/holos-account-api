package api

import (
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/service"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database/pkg/transaction"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler"
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/middleware"
	"github.com/atsumarukun/holos-account-api/internal/app/api/usecase"
)

var (
	healthHdl        handler.HealthHandler
	accountHdl       handler.AccountHandler
	sessionHdl       handler.SessionHandler
	authenticationMW middleware.AuthenticationMiddleware
)

func inject(db *sqlx.DB) {
	transactionObj := transaction.NewDBTransactionObject(db)

	healthHdl = handler.NewHealthHandler()

	accountRepo := database.NewDBAccountRepository(db)
	accountServ := service.NewAccountService(accountRepo)
	accountUC := usecase.NewAccountUsecase(transactionObj, accountRepo, accountServ)
	accountHdl = handler.NewAccountHandler(accountUC)

	sessionRepo := database.NewDBSessionRepository(db)
	sessionUC := usecase.NewSessionUsecase(transactionObj, sessionRepo, accountRepo)
	sessionHdl = handler.NewSessionHandler(sessionUC)

	authenticationMW = middleware.NewAuthenticationMiddleware(sessionUC)
}
