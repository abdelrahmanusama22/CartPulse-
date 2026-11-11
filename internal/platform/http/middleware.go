package httpserver

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var b [16]byte
		if _, err := rand.Read(b[:]); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		id := hex.EncodeToString(b[:])
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func requestLog(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		// Do not log query strings, headers, bodies, cookies or connection URLs.
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		log.Info("http_request", "request_id", c.GetString("request_id"),
			"method", c.Request.Method, "route", route,
			"status", c.Writer.Status(), "duration_ms", time.Since(start).Milliseconds())
	}
}

func recoverPanics(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				log.Error("request_panicked", "request_id", c.GetString("request_id"))
				if !c.Writer.Written() {
					WriteError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred.")
				} else {
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}

func bodyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	}
}
