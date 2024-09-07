package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "authorization"
)

func AttachCORS(r *gin.Engine, allowedOrigins []string, exposedHeaders []string) {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = allowedOrigins
	cfg.AllowCredentials = true
	cfg.AllowHeaders = append(cfg.AllowHeaders, authorizationHeader)
	cfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	for _, header := range exposedHeaders {
		cfg.ExposeHeaders = append(cfg.ExposeHeaders, header)
	}

	r.Use(cors.New(cfg))
}
