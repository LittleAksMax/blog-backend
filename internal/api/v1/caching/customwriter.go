package caching

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type cacheWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w cacheWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
