package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	r.GET("/health", healthHdl.Health)
	r.POST("/login", sessionHdl.Login)
	r.DELETE("/logout", authenticationMW.Authenticate, sessionHdl.Logout)
	r.GET("/authorization", authenticationMW.Authenticate, sessionHdl.Authorize)

	accounts := r.Group("accounts")
	accounts.POST("/", accountHdl.Create)
	accounts.DELETE("/", authenticationMW.Authenticate, accountHdl.Delete)
	accounts.PUT("/name", authenticationMW.Authenticate, accountHdl.UpdateName)
	accounts.PUT("/password", authenticationMW.Authenticate, accountHdl.UpdatePassword)
}
