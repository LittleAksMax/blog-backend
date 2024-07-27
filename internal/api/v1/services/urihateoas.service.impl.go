package services

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type UriServiceImpl struct {
}

func NewUriServiceImpl() *UriServiceImpl {
	return &UriServiceImpl{}
}

func (us *UriServiceImpl) GetPrevUri(ctx *gin.Context, pageNum int, pageSize int) (string, bool) {
	if pageNum == 1 {
		return "", false
	}

	uri := generateUri(ctx.Request.URL.Host+ctx.Request.URL.Path, pageSize, pageNum-1, ctx.Request.URL.Query())
	return uri, true
}

func (us *UriServiceImpl) GetNextUri(ctx *gin.Context, pageNum int, pageSize int, totalCount int) (string, bool) {
	if pageNum*pageSize >= totalCount {
		return "", false
	}

	uri := generateUri(ctx.Request.URL.Host+ctx.Request.URL.Path, pageSize, pageNum+1, ctx.Request.URL.Query())
	return uri, true
}

// NOTE: `base` parameter should exclude the query operator '?'
func generateUri(base string, pageSize int, pageNum int, query map[string][]string) string {
	// manually add page_size and page_num
	uri := base + "?page_size=" + strconv.Itoa(pageSize) + "&page_num=" + strconv.Itoa(pageNum)

	// iterate through query and add all non-pagination queries
	for k, vals := range query {
		if k != "page_size" && k != "page_num" {
			// add all existing values
			// NOTE: we can safely prefix '&' because we manually do
			//       the page_size and page_num parameters
			for _, val := range vals {
				uri += "&" + k + "=" + val
			}
		}
	}
	return uri
}
