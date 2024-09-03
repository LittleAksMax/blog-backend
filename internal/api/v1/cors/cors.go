package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "authorization"
)

func AttachCORS(r *gin.Engine, allowedOrigins []string) {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = allowedOrigins
	cfg.AllowCredentials = true
	cfg.AllowHeaders = append(cfg.AllowHeaders, authorizationHeader)

	r.Use(cors.New(cfg))
}
