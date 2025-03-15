package api

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler"
	"github.com/jmoiron/sqlx"
)

var healthHandler handler.HealthHandler

func inject(_db *sqlx.DB) {
	healthHandler = handler.NewHealthHandler()
}
