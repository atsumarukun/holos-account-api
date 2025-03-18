package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	health := r.Group("health")
	health.GET("/", healthHdl.Health)

	accounts := r.Group("accounts")
	accounts.POST("/", accountHdl.Create)
}
