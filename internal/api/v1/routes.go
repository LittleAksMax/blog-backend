package v1

import (
	"firebase.google.com/go/v4/auth"
	authMiddleware "github.com/LittleAksMax/blog-backend/internal/api/auth"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/caching"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/controllers"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/validators"
	fbAuth "github.com/LittleAksMax/blog-backend/internal/auth"
	"github.com/LittleAksMax/blog-backend/internal/cache"
	"github.com/gin-gonic/gin"
	"time"
)

func addPostsRoutes(versionGroup *gin.RouterGroup, pc *controllers.PostController, cm *caching.CacheManager, authClient *auth.Client) {
	postsGroup := versionGroup.Group("/posts")
	{
		postsGroup.GET("/", authMiddleware.ReadToken(authClient), authMiddleware.ReadAdmin, validators.QueryValidate, cm.Cache(time.Minute*5, caching.HashGetAllPosts), pc.GetPosts)
		postsGroup.GET("/:idOrSlug", authMiddleware.ReadToken(authClient), authMiddleware.ReadAdmin, validators.RouteIdOrSlugValidate(), cm.Cache(time.Minute*5, caching.HashGetPost), pc.GetPost)
		postsGroup.POST("/", authMiddleware.ReadToken(authClient), authMiddleware.ReadAdmin, authMiddleware.RequiresAdmin, validators.ReqValidate[*models.CreatePostRequest], pc.CreatePost)
		postsGroup.PUT("/:id", authMiddleware.ReadToken(authClient), authMiddleware.ReadAdmin, authMiddleware.RequiresAdmin, validators.RouteIdValidate(), validators.ReqValidate[*models.UpdatePostRequest], validators.UpdateReqValidate, pc.UpdatePost)
		postsGroup.PUT("/:id/archive", authMiddleware.ReadToken(authClient), authMiddleware.ReadAdmin, authMiddleware.RequiresAdmin, validators.RouteIdValidate(), pc.ArchivePost)
		postsGroup.PUT("/:id/publish", authMiddleware.ReadToken(authClient), authMiddleware.ReadAdmin, authMiddleware.RequiresAdmin, validators.RouteIdValidate(), pc.PublishPost)
		postsGroup.DELETE("/:id", authMiddleware.ReadToken(authClient), authMiddleware.ReadAdmin, authMiddleware.RequiresAdmin, validators.RouteIdValidate(), pc.DeletePost)
	}
}

func AttachVersion(api *gin.RouterGroup, pc *controllers.PostController, cacheCfg *cache.Config, authCfg *fbAuth.Config) {
	// create the cache manager used
	cm := caching.NewCacheManager(cacheCfg)
	authClient := authCfg.AuthClient

	versionGroup := api.Group("/v1")
	{
		addPostsRoutes(versionGroup, pc, cm, authClient)
	}
}
