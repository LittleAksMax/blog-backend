package health

import (
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type HealthController struct {
	db *mongo.Database
	// TODO: redis field
}

func NewHealthController(dbCfg *db.Config) *HealthController {
	return &HealthController{db: dbCfg.Database}
	// TODO: set up redis field
}

func (hc *HealthController) Health(ctx *gin.Context) {
	res := gin.H{"status": "ok"}
	dbStatus, dbOk := checkDbHealth(ctx.Request.Context(), hc.db)
	cacheStatus, cacheOk := checkCacheHealth(ctx.Request.Context()) // TODO: redis parameter

	status := http.StatusOK
	if !dbOk || !cacheOk {
		status = http.StatusServiceUnavailable
		res["status"] = "not ok"
	}
	res["db"] = dbStatus
	res["cache"] = cacheStatus

	ctx.JSON(status, res)
}
