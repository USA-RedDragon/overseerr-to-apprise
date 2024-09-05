package server

import (
	"log/slog"
	"net/http"

	"github.com/USA-RedDragon/overseerr-to-apprise/internal/config"
	"github.com/gin-gonic/gin"
)

func applyRoutes(r *gin.Engine, config *config.Config) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	apiV1 := r.Group("/v1")
	v1(apiV1, config)

	r.NoRoute(func(c *gin.Context) {
		slog.Warn("Not Found", "path", c.Request.URL.Path)
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
}

func v1(group *gin.RouterGroup, config *config.Config) {
}
