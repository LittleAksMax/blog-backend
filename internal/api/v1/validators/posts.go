package validators

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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
	ctx.Set("request", req)
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
