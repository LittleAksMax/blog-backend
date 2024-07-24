package controllers

import (
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/services"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type PostController struct {
	ps  services.PostService
	cps services.ContentParserService
}

func NewPostController(ps services.PostService, cps services.ContentParserService) *PostController {
	return &PostController{ps: ps, cps: cps}
}

func (pc *PostController) GetPosts(ctx *gin.Context) {
	pf := models.PagedQuery{} // TODO: utility for creating this from validation middleware
	pq := models.PostQuery{}  // TODO: utility for creating this from validation middleware

	posts, err := pc.ps.GetPosts(ctx.Request.Context(), &pf, &pq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, posts)
}

func (pc *PostController) GetPost(ctx *gin.Context) {
	idParam, exists := ctx.Get("id")
	if !exists {
		panic("parameter 'id' must be set")
	}
	id := idParam.(primitive.ObjectID)

	post, err := pc.ps.GetPost(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post with ID " + id.Hex() + " not found."})
	}

	ctx.JSON(http.StatusOK, post)
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	request, exists := ctx.Get("request")
	if !exists {
		panic("parameter 'request' must be set")
	}
	cpr := request.(*models.CreatePostRequest)

	currentDate := time.Now()
	// add to mongo
	dto := models.PostDto{
		Title:       cpr.Title,
		Content:     cpr.Content,
		Media:       pc.cps.GetMediaLinks(cpr.Content),
		Collections: cpr.Collections,
		Tags:        cpr.Tags,
		LastUpdated: currentDate,
		Status:      db.Published.String(), // assuming anything that gets updated is to be published
		Featured:    *cpr.Featured,
	}

	err := pc.ps.CreatePost(ctx.Request.Context(), &dto)

	if err != nil {
		// TODO: handle error properly
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// The work completed without cancellation
	ctx.JSON(http.StatusOK, dto)
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	idParam, exists := ctx.Get("id")
	if !exists {
		panic("parameter 'id' must be set")
	}
	id := idParam.(primitive.ObjectID)

	request, exists := ctx.Get("request")
	if !exists {
		panic("parameter 'request' must be set")
	}
	upr := request.(*models.UpdatePostRequest)

	currentDate := time.Now()
	// add to mongo
	dto := models.PostDto{
		Title:       upr.Title,
		Content:     upr.Content,
		Media:       pc.cps.GetMediaLinks(upr.Content),
		Collections: upr.Collections,
		Tags:        upr.Tags,
		Published:   currentDate,
		LastUpdated: currentDate,
		Status:      db.Published.String(),
		Featured:    *upr.Featured,
	}

	err := pc.ps.UpdatePost(ctx.Request.Context(), id, &dto)

	if err != nil {
		// TODO: handle error properly
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// The work completed without cancellation
	ctx.JSON(http.StatusOK, dto)
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	idParam, exists := ctx.Get("id")
	if !exists {
		panic("parameter 'id' must be set")
	}
	id := idParam.(primitive.ObjectID)

	err := pc.ps.DeletePost(ctx, id)

	if err != nil {
		// TODO: handle error
		// TODO: handle `not found` and other errors respectively
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.Status(http.StatusNoContent)
}
