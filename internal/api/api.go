package api

import (
	"fmt"
	"github.com/LittleAksMax/blog-backend/internal/api/health"
	v1 "github.com/LittleAksMax/blog-backend/internal/api/v1"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/controllers"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/services"
	fbAuth "github.com/LittleAksMax/blog-backend/internal/auth"
	"github.com/LittleAksMax/blog-backend/internal/cache"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func RunApi(port int, apiKey string, dbCfg *db.Config, cacheCfg *cache.Config, authCfg *fbAuth.Config) {
	r := gin.Default()

	// create all relevant controllers and services for API
	pc, hc := createControllers(dbCfg, cacheCfg)

	// configure manual health checks
	healthGroup := r.Group("/")
	{
		health.AttachHealthChecks(healthGroup, hc, apiKey)
	}

	// configure routes using controllers
	apiGroup := r.Group("/api")
	{
		v1.AttachVersion(apiGroup, pc, apiKey, cacheCfg, authCfg)
	}

	addr := fmt.Sprintf(":%d", port)
	err := r.Run(addr)

	if err != nil {
		panic(fmt.Sprintf("Failed to start server on address %s", addr))
	}
}

func createControllers(dbCfg *db.Config, cacheCfg *cache.Config) (*controllers.PostController, *health.HealthController) {
	ps := services.NewPostServiceImpl(dbCfg)
	cps := services.NewContentParserServiceImpl()
	us := services.NewUriServiceImpl()
	pc := controllers.NewPostController(ps, cps, us)
	hc := health.NewHealthController(dbCfg, cacheCfg)
	return pc, hc
}
