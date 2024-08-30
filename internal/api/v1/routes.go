package v1

import (
	authMiddleware "github.com/LittleAksMax/blog-backend/internal/api/auth"
	caching2 "github.com/LittleAksMax/blog-backend/internal/api/v1/middleware/caching"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/controllers"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/middleware"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	fbAuth "github.com/LittleAksMax/blog-backend/internal/auth"
	"github.com/LittleAksMax/blog-backend/internal/cache"
	"github.com/gin-gonic/gin"
)

func addPostsRoutes(versionGroup *gin.RouterGroup, pc *controllers.PostController, cm *caching2.CacheManager, authClient *auth.Client, apiKey string) {
	postsGroup := versionGroup.Group("/posts")
	{
		postsGroup.GET("/", authMiddleware.RequiresAPIKey(apiKey), middleware.QueryValidate, cm.Cache(time.Minute*1, caching2.HashGetAllPosts), pc.GetPosts)
		postsGroup.GET("/:idOrSlug", authMiddleware.RequiresAPIKey(apiKey), middleware.RouteIdOrSlugValidate(), cm.Cache(time.Minute*1, caching2.HashGetPost), pc.GetPost)
		postsGroup.POST("/", authMiddleware.RequiresToken(authClient), authMiddleware.RequiresAdmin, middleware.ReqValidate[*models.CreatePostRequest], pc.CreatePost)
		postsGroup.PUT("/:id", authMiddleware.RequiresToken(authClient), authMiddleware.RequiresAdmin, middleware.RouteIdValidate(), middleware.ReqValidate[*models.UpdatePostRequest], pc.UpdatePost)
		postsGroup.DELETE("/:id", authMiddleware.RequiresToken(authClient), authMiddleware.RequiresAdmin, middleware.RouteIdValidate(), pc.DeletePost)
	}
}

func AttachVersion(api *gin.RouterGroup, pc *controllers.PostController, apiKey string, cacheCfg *cache.Config, authCfg *fbAuth.Config) {
	// create the cache manager used
	cm := caching2.NewCacheManager(cacheCfg)
	authClient := authCfg.AuthClient

	versionGroup := api.Group("/v1")
	{
		addPostsRoutes(versionGroup, pc, cm, authClient, apiKey)
	}
}
