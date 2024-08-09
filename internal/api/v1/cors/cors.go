package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func AttachCORS(r *gin.Engine, allowedOrigins []string) {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = allowedOrigins

	r.Use(cors.New(cfg))
}
