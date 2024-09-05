package server

import (
	"log/slog"

	"github.com/USA-RedDragon/overseerr-to-apprise/internal/config"
	"github.com/USA-RedDragon/overseerr-to-apprise/internal/metrics"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func applyMiddleware(
	r *gin.Engine,
	config *config.Config,
	otelComponent string,
	metrics *metrics.Metrics) {
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/health", "/metrics"}}))
	r.TrustedPlatform = "X-Real-IP"

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "authorization")
	corsConfig.AllowCredentials = true
	corsConfig.AllowWildcard = true
	if len(config.HTTP.CORSHosts) == 0 {
		corsConfig.AllowAllOrigins = true
	}
	corsConfig.AllowOrigins = config.HTTP.CORSHosts
	r.Use(cors.New(corsConfig))

	err := r.SetTrustedProxies(config.HTTP.TrustedProxies)
	if err != nil {
		slog.Error("Failed to set trusted proxies", "error", err.Error())
	}

	r.Use(providerMiddleware("config", config))
	r.Use(providerMiddleware("metrics", metrics))

	if config.HTTP.Tracing.Enabled {
		r.Use(otelgin.Middleware(otelComponent))
		r.Use(tracingProvider(config))
	}
}

func providerMiddleware[T any](key string, toProvide T) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, toProvide)
		c.Next()
	}
}

func tracingProvider(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.HTTP.Tracing.OTLPEndpoint != "" {
			ctx := c.Request.Context()
			span := trace.SpanFromContext(ctx)
			if span.IsRecording() {
				span.SetAttributes(
					attribute.String("http.method", c.Request.Method),
					attribute.String("http.path", c.Request.URL.Path),
				)
			}
		}
		c.Next()
	}
}
