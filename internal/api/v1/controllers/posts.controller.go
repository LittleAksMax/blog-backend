package controllers

import (
	"errors"
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
	us  services.UriService
}

func NewPostController(ps services.PostService, cps services.ContentParserService, us services.UriService) *PostController {
	return &PostController{ps: ps, cps: cps, us: us}
}

func (pc *PostController) GetPosts(ctx *gin.Context) {
	pf := ctx.MustGet("pf").(*models.PostFilter) // post filter
	pq := ctx.MustGet("pq").(*models.PagedQuery) // pagination query

	posts, totalCount, err := pc.ps.GetPosts(ctx.Request.Context(), pq, pf)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// form next and previous page pointers
	var prevUri, nextUri any
	if prev, ok := pc.us.GetPrevUri(ctx, pq.PageNum, pq.PageSize); !ok {
		prevUri = nil
	} else {
		prevUri = prev
	}
	if next, ok := pc.us.GetNextUri(ctx, pq.PageNum, pq.PageSize, totalCount); !ok {
		nextUri = nil
	} else {
		nextUri = next
	}
	ctx.JSON(http.StatusOK, gin.H{"data": posts, "prev": prevUri, "next": nextUri})
}

func (pc *PostController) GetPost(ctx *gin.Context) {
	id := ctx.MustGet("id").(primitive.ObjectID)

	post, err := pc.ps.GetPost(ctx, id)

	if err != nil {
		var nfErr services.NotFoundErr
		if errors.Is(err, &nfErr) {
			// changes are the object was not found in the database
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	cpr := ctx.MustGet("request").(*models.CreatePostRequest)

	currentDate := time.Now()
	// add to mongo
	dto := models.PostDto{
		Title:       cpr.Title,
		Content:     cpr.Content,
		Media:       pc.cps.GetMediaLinks(cpr.Content),
		Collections: cpr.Collections,
		Tags:        cpr.Tags,
		LastUpdated: currentDate,
		Published:   currentDate,
		Status:      db.Published.String(), // assuming anything that gets updated is to be published
		Featured:    *cpr.Featured,
	}

	err := pc.ps.CreatePost(ctx.Request.Context(), &dto)

	if err != nil {
		// chances are that caught error is a conflict error
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// The work completed without cancellation
	ctx.JSON(http.StatusOK, dto)
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	id := ctx.MustGet("id").(primitive.ObjectID)
	upr := ctx.MustGet("request").(*models.UpdatePostRequest)

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
		var nfErr services.NotFoundErr
		if errors.Is(err, &nfErr) {
			// changes are the object was not found in the database
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	// The work completed without cancellation
	ctx.JSON(http.StatusOK, dto)
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	id := ctx.MustGet("id").(primitive.ObjectID)

	err := pc.ps.DeletePost(ctx, id)

	if err != nil {
		var nfErr services.NotFoundErr
		if errors.Is(err, &nfErr) {
			// changes are the object was not found in the database
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
