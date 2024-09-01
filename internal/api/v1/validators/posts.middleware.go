package validators

import (
	"github.com/LittleAksMax/blog-backend/internal/api/v1/models"
	"github.com/LittleAksMax/blog-backend/internal/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

const (
	RequestKey  = "request"
	IdOrSlugKey = "idOrSlug"
)

func ReqValidate[T any](ctx *gin.Context) {
	var req T
	diag, ok := validateRequest(ctx, &req)

	// return if validation failed
	if !ok {
		var res gin.H
		single, ok := diag["_"]
		if ok {
			res = gin.H{"error": single}
		} else {
			res = gin.H{"errors": diag}
		}
		ctx.JSON(http.StatusBadRequest, res)
		ctx.Abort()
		return
	}

	// add to context and call next
	ctx.Set(RequestKey, req)
	ctx.Next()
}

func UpdateReqValidate(ctx *gin.Context) {
	upr := ctx.MustGet(RequestKey).(*models.UpdatePostRequest)

	_, ok := db.PostStatusFromString(upr.Status)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": upr.Status + " is not a valid post status."})
		return
	}

	ctx.Next()
}

func RouteIdParamValidateWithParam(idParamName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param(idParamName)
		id, err := primitive.ObjectIDFromHex(idParam)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		// add to context and call next
		ctx.Set(idParamName, id)
		ctx.Next()
	}
}

func RouteIdValidate() gin.HandlerFunc {
	return RouteIdParamValidateWithParam("id")
}

func RouteIdOrSlugValidateWithParamName(idOrSlugParamName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idOrSlug := ctx.Param(idOrSlugParamName)
		id, err := primitive.ObjectIDFromHex(idOrSlug)

		// set the slug if what we got is not a valid ID, otherwise set the id
		if err != nil {
			ctx.Set(IdOrSlugKey, false) // we should use the slug
			ctx.Set("slug", idOrSlug)
		} else {
			ctx.Set(IdOrSlugKey, true) // we should use the ID
			ctx.Set("id", id)
		}

		ctx.Next()
	}
}

func RouteIdOrSlugValidate() gin.HandlerFunc {
	return RouteIdOrSlugValidateWithParamName("idOrSlug")
}

func QueryValidate(ctx *gin.Context) {
	pq := models.PagedQuery{
		PageNum:  models.DefaultPageNumber,
		PageSize: models.DefaultPageSize,
	}
	pf := models.PostFilter{
		Title:       nil,
		Tags:        nil,
		Collections: nil,
		Featured:    nil,
	}

	if err := ctx.BindQuery(&pf); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if err := ctx.BindQuery(&pq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Set("pf", &pf)
	ctx.Set("pq", &pq)
	ctx.Next()
}
