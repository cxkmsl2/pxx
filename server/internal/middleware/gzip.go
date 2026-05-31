package middleware

import (
	"compress/gzip"
	"io"
		"strings"

	"github.com/gin-gonic/gin"
)

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}
		gw := gzip.NewWriter(c.Writer)
		c.Writer = &gzipWriter{ResponseWriter: c.Writer, Writer: gw}
		c.Header("Content-Encoding", "gzip")
		defer gw.Close()
		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	io.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}