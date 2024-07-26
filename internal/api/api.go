package api

import (
	"fmt"
	"github.com/LittleAksMax/blog-backend/internal/api/health"
	v1 "github.com/LittleAksMax/blog-backend/internal/api/v1"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/controllers"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/services"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func RunApi(port int, dbCfg *db.Config) {
	r := gin.Default()

	// create all relevant controllers and services for API
	pc, hc := createControllers(dbCfg)

	// configure manual health checks
	healthGroup := r.Group("/")
	{
		health.AttachHealthChecks(healthGroup, hc)
	}

	// configure routes using controllers
	apiGroup := r.Group("/api")
	{
		v1.AttachVersion(apiGroup, pc)
	}

	addr := fmt.Sprintf(":%d", port)
	err := r.Run(addr)

	if err != nil {
		panic(fmt.Sprintf("Failed to start server on address %s", addr))
	}
}

func createControllers(dbCfg *db.Config) (*controllers.PostController, *health.HealthController) {
	ps := services.NewPostServiceImpl(dbCfg)
	cps := services.NewContentParserServiceImpl()
	pc := controllers.NewPostController(ps, cps)
	hc := health.NewHealthController(dbCfg)
	return pc, hc
}
