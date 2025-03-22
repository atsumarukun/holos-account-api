package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	r.GET("/health", healthHdl.Health)
	r.POST("/login", sessionHdl.Login)
	r.DELETE("/logout", authenticationMW.Authenticate, sessionHdl.Logout)

	accounts := r.Group("accounts")
	accounts.POST("/", accountHdl.Create)
}
