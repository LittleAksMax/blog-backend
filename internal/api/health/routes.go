package health

import (
	authMiddleware "github.com/LittleAksMax/blog-backend/internal/api/auth"
	"github.com/gin-gonic/gin"
)

func AttachHealthChecks(api *gin.RouterGroup, hc *HealthController, apiKey string) {
	api.GET("/healthz", authMiddleware.RequiresAPIKey(apiKey), hc.Health)
}
