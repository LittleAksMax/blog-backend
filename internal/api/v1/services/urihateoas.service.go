package services

import "github.com/gin-gonic/gin"

type UriService interface {
	GetPrevUri(ctx *gin.Context, pageNum int, pageSize int) (string, bool)
	GetNextUri(ctx *gin.Context, pageNum int, pageSize int, totalCount int) (string, bool)
}
