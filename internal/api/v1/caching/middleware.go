package caching

import (
	"bytes"
	"errors"
	"github.com/LittleAksMax/blog-backend/internal/api/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

type RequestHashFunc func(*gin.Context) string

func (cm *CacheManager) Cache(duration time.Duration, hasher RequestHashFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.GetHeader(auth.BearerTokenKey)

		// we should never cache when there is auth header present
		if auth.CheckExists(bearer) {
			ctx.Next()
			return
		}

		// get request key for checking if stored in cache
		key := hasher(ctx)

		cached, err := cm.rdb.Get(ctx.Request.Context(), key).Result()

		// if no redis.Nil error then key does exist, and we should return the contents
		if !errors.Is(err, redis.Nil) {
			ctx.Data(http.StatusOK, "application/json", []byte(cached))
			ctx.Abort()
			return
		}

		// use custom response writer, so I can pull raw JSON written
		writer := &cacheWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}

		// run the handler
		// NOTE: this might overwrite the previous Writer's data
		ctx.Writer = writer // replace the context's response writer with our custom writer
		ctx.Next()

		// don't store in cache if the status is not OK
		if ctx.Writer.Status() != http.StatusOK {
			return
		}

		// store created status
		contents := writer.body.String()
		err = cm.rdb.Set(ctx.Request.Context(), key, contents, duration).Err()
		if err != nil {
			log.Println("Cache failed: " + err.Error())
		}
	}
}

func (cm *CacheManager) InvalidateCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		cm.rdb.Del(ctx.Request.Context())
		// TODO: figure out adding caches with a tag
		// TODO: evict cached by tag
	}
}
