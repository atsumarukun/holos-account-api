package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	r.GET("/health", healthHdl.Health)

	accounts := r.Group("accounts")
	accounts.POST("/", accountHdl.Create)
	accounts.DELETE("/", authenticationMW.Authenticate, accountHdl.Delete)
	accounts.PATCH("/name", authenticationMW.Authenticate, accountHdl.UpdateName)
	accounts.PATCH("/password", authenticationMW.Authenticate, accountHdl.UpdatePassword)

	sessions := r.Group("sessions")
	sessions.POST("/", sessionHdl.Create)
	sessions.DELETE("/", authenticationMW.Authenticate, sessionHdl.Delete)
	sessions.GET("/verify", authenticationMW.Authenticate, sessionHdl.Verify)
}
