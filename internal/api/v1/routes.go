package v1

import (
	"github.com/LittleAksMax/blog-backend/internal/api/v1/controllers"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/validators"
	"github.com/gin-gonic/gin"
)

func addPostsRoutes(versionGroup *gin.RouterGroup, pc *controllers.PostController) {
	postsGroup := versionGroup.Group("/posts")
	{
		postsGroup.GET("/", pc.GetPosts)
		postsGroup.GET("/:id", validators.RouteIdValidate(), pc.GetPost)
		postsGroup.POST("/", validators.ReqValidate[*models.CreatePostRequest], pc.CreatePost)
		postsGroup.PUT("/:id", validators.RouteIdValidate(), validators.ReqValidate[*models.UpdatePostRequest], pc.UpdatePost)
		postsGroup.DELETE("/:id", validators.RouteIdValidate(), pc.DeletePost)
	}
}

func AttachVersion(api *gin.RouterGroup, pc *controllers.PostController) {
	versionGroup := api.Group("/v1")
	{
		addPostsRoutes(versionGroup, pc)
	}
}
