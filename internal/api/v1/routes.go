package v1

import (
	"github.com/LittleAksMax/blog-backend/internal/api/caching"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/controllers"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/validators"
	"github.com/LittleAksMax/blog-backend/internal/cache"
	"github.com/gin-gonic/gin"
	"time"
)

func addPostsRoutes(versionGroup *gin.RouterGroup, pc *controllers.PostController, cm *caching.CacheManager) {
	postsGroup := versionGroup.Group("/posts")
	{
		postsGroup.GET("/", validators.QueryValidate, cm.Cache(time.Minute*1), pc.GetPosts)
		postsGroup.GET("/:id", validators.RouteIdValidate(), cm.Cache(time.Minute*1), pc.GetPost)
		postsGroup.POST("/", validators.ReqValidate[*models.CreatePostRequest], pc.CreatePost)
		postsGroup.PUT("/:id", validators.RouteIdValidate(), validators.ReqValidate[*models.UpdatePostRequest], pc.UpdatePost)
		postsGroup.DELETE("/:id", validators.RouteIdValidate(), pc.DeletePost)
	}
}

func AttachVersion(api *gin.RouterGroup, pc *controllers.PostController, cacheCfg *cache.Config) {
	// create the cache manager used
	cm := caching.NewCacheManager(cacheCfg)

	versionGroup := api.Group("/v1")
	{
		addPostsRoutes(versionGroup, pc, cm)
	}
}
