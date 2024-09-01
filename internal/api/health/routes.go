package health

import (
	"github.com/gin-gonic/gin"
)

func AttachHealthChecks(api *gin.RouterGroup, hc *HealthController) {
	api.GET("/healthz", hc.Health)
}
