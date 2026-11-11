// Package httpserver contains transport conventions, not business logic.
package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReadinessCheck func(context.Context) error

func NewRouter(log *slog.Logger, ready ReadinessCheck) *gin.Engine {
	r := gin.New()
	// Only explicitly configured reverse proxies may be trusted in a deployment.
	_ = r.SetTrustedProxies(nil)
	r.HandleMethodNotAllowed = true
	r.Use(requestID(), requestLog(log), recoverPanics(log), bodyLimit())
	r.NoRoute(func(c *gin.Context) {
		WriteError(c, http.StatusNotFound, "NOT_FOUND", "The requested resource was not found.")
	})
	r.NoMethod(func(c *gin.Context) {
		WriteError(c, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "This method is not allowed.")
	})
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	r.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if ready == nil || ready(ctx) != nil {
			WriteError(c, http.StatusServiceUnavailable, "NOT_READY", "A required dependency is unavailable.")
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
	// Business routes are deliberately not mounted until their modules implement
	// validation, services, ownership checks and the correct auth middleware.
	// See api/README.md for the agreed /api/v1 route groups.
	return r
}
