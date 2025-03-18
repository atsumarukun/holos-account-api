package api

import (
	"github.com/jmoiron/sqlx"

	"github.com/atsumarukun/holos-account-api/internal/app/api/interface/handler"
)

var healthHandler handler.HealthHandler

func inject(_db *sqlx.DB) {
	healthHandler = handler.NewHealthHandler()
}
