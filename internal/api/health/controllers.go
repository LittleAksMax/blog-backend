package health

import (
	"github.com/LittleAksMax/blog-backend/internal/cache"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type HealthController struct {
	db  *mongo.Database
	rdb *redis.Client
}

func NewHealthController(dbCfg *db.Config, cacheCfg *cache.Config) *HealthController {
	return &HealthController{db: dbCfg.Database, rdb: cacheCfg.Client}
}

func (hc *HealthController) Health(ctx *gin.Context) {
	res := gin.H{"status": "ok"}
	dbStatus, dbOk := checkDbHealth(ctx.Request.Context(), hc.db)
	cacheStatus, cacheOk := checkCacheHealth(ctx.Request.Context(), hc.rdb)

	status := http.StatusOK
	if !dbOk || !cacheOk {
		status = http.StatusServiceUnavailable
		res["status"] = "not ok"
	}
	res["db"] = dbStatus
	res["cache"] = cacheStatus

	ctx.JSON(status, res)
}
