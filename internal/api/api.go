package api

import (
	"fmt"
	v1 "github.com/LittleAksMax/blog-backend/internal/api/v1"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/controllers"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/services"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func RunApi(port int, dbCfg *db.Config) {
	r := gin.Default()

	// create all relevant controllers and services
	pc := createControllers(dbCfg)

	// configure routes using controllers
	g := r.Group("/api")
	{
		v1.AttachVersion(g, pc)
	}

	addr := fmt.Sprintf(":%d", port)
	err := r.Run(addr)

	if err != nil {
		panic(fmt.Sprintf("Failed to start server on address %s", addr))
	}
}

func createControllers(dbCfg *db.Config) *controllers.PostController {
	ps := services.NewPostServiceImpl(dbCfg)
	cps := services.NewContentParserServiceImpl()
	pc := controllers.NewPostController(ps, cps)
	return pc
}
