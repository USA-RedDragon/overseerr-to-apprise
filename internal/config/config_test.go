package config_test

import (
	"context"
	"errors"
	"testing"

	"github.com/USA-RedDragon/overseerr-to-apprise/cmd"
	"github.com/USA-RedDragon/overseerr-to-apprise/internal/config"
)

//nolint:golint,gochecknoglobals
var requiredFlags = []string{}

func TestExampleConfig(t *testing.T) {
	t.Parallel()
	cmd := cmd.NewCommand("testing", "deadbeef")
	cmd.SetContext(context.Background())
	err := cmd.ParseFlags([]string{"--config", "../../config.example.yaml", "--log_level", "debug"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err := config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TesMissingOLTPEndpoint(t *testing.T) {
	t.Parallel()

	cmd := cmd.NewCommand("testing", "deadbeef")
	cmd.SetContext(context.Background())
	err := cmd.ParseFlags(append([]string{"--http.tracing.enabled", "true"}, requiredFlags...))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err := config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); !errors.Is(err, config.ErrOTLPEndpointRequired) {
		t.Errorf("unexpected error: %v", err)
	}

	err = cmd.ParseFlags(append([]string{"--http.tracing.enabled", "true", "--http.tracing.otlp_endpoint", "dummy"}, requiredFlags...))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	testConfig, err = config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := testConfig.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Parallel tests are not allowed with t.Setenv
//
//nolint:golint,paralleltest
func TestEnvConfig(t *testing.T) {
	cmd := cmd.NewCommand("testing", "deadbeef")
	cmd.SetContext(context.Background())
	t.Setenv("HTTP__PORT", "8087")
	t.Setenv("HTTP__METRICS__PORT", "8088")
	t.Setenv("HTTP__METRICS__IPV4_HOST", "0.0.0.0")
	t.Setenv("HTTP__METRICS__IPV6_HOST", "::0")
	t.Setenv("HTTP__IPV4_HOST", "127.0.0.1")
	t.Setenv("HTTP__IPV6_HOST", "::1")
	t.Setenv("HTTP__PPROF__ENABLED", "true")
	t.Setenv("HTTP__TRUSTED_PROXIES", "127.0.0.1,127.0.0.2")
	t.Setenv("HTTP__METRICS__ENABLED", "true")
	t.Setenv("HTTP__TRACING__ENABLED", "true")
	t.Setenv("HTTP__TRACING__OTLP_ENDPOINT", "http://localhost:4317")
	t.Setenv("HTTP__CORS_HOSTS", "http://localhost:8080,http://localhost:8081")
	t.Setenv("HTTP__BACKEND_URL", "http://localhost:8081")
	t.Setenv("LOG_LEVEL", "warn")

	config, err := config.LoadConfig(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if config.HTTP.Port != 8087 {
		t.Errorf("unexpected HTTP port: %d", config.HTTP.Port)
	}
	if config.HTTP.Metrics.Port != 8088 {
		t.Errorf("unexpected HTTP metrics port: %d", config.HTTP.Metrics.Port)
	}
	if config.HTTP.Metrics.IPV4Host != "0.0.0.0" {
		t.Errorf("unexpected HTTP metrics IPv4 host: %s", config.HTTP.Metrics.IPV4Host)
	}
	if config.HTTP.Metrics.IPV6Host != "::0" {
		t.Errorf("unexpected HTTP metrics IPv6 host: %s", config.HTTP.Metrics.IPV6Host)
	}
	if config.HTTP.IPV4Host != "127.0.0.1" {
		t.Errorf("unexpected HTTP IPv4 host: %s", config.HTTP.IPV4Host)
	}
	if config.HTTP.IPV6Host != "::1" {
		t.Errorf("unexpected HTTP IPv6 host: %s", config.HTTP.IPV6Host)
	}
	if !config.HTTP.PProf.Enabled {
		t.Error("unexpected HTTP pprof enabled")
	}
	if len(config.HTTP.TrustedProxies) != 2 {
		t.Errorf("unexpected HTTP trusted proxies: %v", config.HTTP.TrustedProxies)
	}
	if config.HTTP.TrustedProxies[0] != "127.0.0.1" {
		t.Errorf("unexpected HTTP trusted proxy: %s", config.HTTP.TrustedProxies[0])
	}
	if config.HTTP.TrustedProxies[1] != "127.0.0.2" {
		t.Errorf("unexpected HTTP trusted proxy: %s", config.HTTP.TrustedProxies[1])
	}
	if !config.HTTP.Metrics.Enabled {
		t.Error("unexpected HTTP metrics enabled")
	}
	if !config.HTTP.Tracing.Enabled {
		t.Error("unexpected HTTP tracing enabled")
	}
	if config.HTTP.Tracing.OTLPEndpoint != "http://localhost:4317" {
		t.Errorf("unexpected HTTP tracing OTLP endpoint: %s", config.HTTP.Tracing.OTLPEndpoint)
	}
	if len(config.HTTP.CORSHosts) != 2 {
		t.Errorf("unexpected HTTP CORS hosts: %v", config.HTTP.CORSHosts)
	}
	if config.HTTP.CORSHosts[0] != "http://localhost:8080" {
		t.Errorf("unexpected HTTP CORS host: %s", config.HTTP.CORSHosts[0])
	}
	if config.HTTP.CORSHosts[1] != "http://localhost:8081" {
		t.Errorf("unexpected HTTP CORS host: %s", config.HTTP.CORSHosts[1])
	}
	if config.LogLevel != "warn" {
		t.Errorf("unexpected log level: %s", config.LogLevel)
	}
}
