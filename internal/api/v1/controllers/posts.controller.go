package controllers

import (
	"errors"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/services"
	"github.com/LittleAksMax/blog-backend/internal/api/v1/validators"
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
	admin := ctx.MustGet("admin").(bool)

	posts, totalCount, err := pc.ps.GetPosts(ctx.Request.Context(), pq, pf, admin)

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
	idUsed := ctx.MustGet(validators.IdOrSlugKey).(bool)
	admin := ctx.MustGet("admin").(bool)

	var post *models.PostDto
	var err error
	if idUsed {
		id := ctx.MustGet("id").(primitive.ObjectID)
		post, err = pc.ps.GetPostById(ctx.Request.Context(), id, admin)
	} else {
		slug := ctx.MustGet("slug").(string)
		post, err = pc.ps.GetPostBySlug(ctx.Request.Context(), slug, admin)
	}

	if err != nil {
		//var nfErr services.NotFoundErr
		//var snfErr services.SlugNotFoundErr
		if errors.Is(err, &services.NotFoundErr{}) || errors.Is(err, &services.SlugNotFoundErr{}) {
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
	cpr := ctx.MustGet(validators.RequestKey).(*models.CreatePostRequest)

	currentDate := time.Now()
	// add to mongo
	dto := models.PostDto{
		Title:       cpr.Title,
		Slug:        models.GenerateSlug(cpr.Title),
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
		var cfErr services.ConflictErr
		if errors.As(err, &cfErr) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// The work completed without cancellation
	ctx.JSON(http.StatusOK, dto)
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	id := ctx.MustGet("id").(primitive.ObjectID)
	upr := ctx.MustGet(validators.RequestKey).(*models.UpdatePostRequest)

	// NOTE: technically, to be in this controller, one must already be an admin
	//       because of the middleware, but this is just for clarity
	admin := ctx.MustGet("admin").(bool)

	// fetch old equivalent (if it exists) to set the published date
	oldDto, err := pc.ps.GetPostById(ctx.Request.Context(), id, admin)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// add to mongo
	currentDate := time.Now()
	dto := models.PostDto{
		Title:       upr.Title,
		Slug:        models.GenerateSlug(upr.Title),
		Content:     upr.Content,
		Media:       pc.cps.GetMediaLinks(upr.Content),
		Collections: upr.Collections,
		Tags:        upr.Tags,
		Published:   oldDto.Published,
		LastUpdated: currentDate,
		Status:      upr.Status,
		Featured:    *upr.Featured,
	}

	err = pc.ps.UpdatePost(ctx.Request.Context(), id, &dto)

	if err != nil {
		var nfErr services.NotFoundErr
		var cfErr services.ConflictErr
		if errors.As(err, &nfErr) {
			// changes are the object was not found in the database
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &cfErr) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// The work completed without cancellation
	ctx.Status(http.StatusOK)
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
