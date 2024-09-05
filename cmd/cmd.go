package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"syscall"
	"time"

	"github.com/USA-RedDragon/overseerr-to-apprise/internal/config"
	"github.com/USA-RedDragon/overseerr-to-apprise/internal/metrics"
	"github.com/USA-RedDragon/overseerr-to-apprise/internal/server"
	"github.com/spf13/cobra"
	"github.com/ztrue/shutdown"
	"golang.org/x/sync/errgroup"
)

func NewCommand(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "overseerr-to-apprise",
		Version: fmt.Sprintf("%s - %s", version, commit),
		Annotations: map[string]string{
			"version": version,
			"commit":  commit,
		},
		RunE:          run,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	config.RegisterFlags(cmd)
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	slog.Info("overseerr-to-apprise", "version", cmd.Annotations["version"], "commit", cmd.Annotations["commit"])

	cfg, err := config.LoadConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch cfg.LogLevel {
	case config.LogLevelDebug:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case config.LogLevelInfo:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case config.LogLevelWarn:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case config.LogLevelError:
		slog.SetLogLoggerLevel(slog.LevelError)
	}

	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	metrics := metrics.NewMetrics()

	slog.Info("Starting HTTP server")
	server := server.NewServer(cfg, metrics)
	err = server.Start()
	if err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	stop := func(_ os.Signal) {
		slog.Info("Shutting down")

		errGrp := errgroup.Group{}

		errGrp.Go(func() error {
			return server.Stop()
		})

		err := errGrp.Wait()
		if err != nil {
			slog.Error("Shutdown error", "error", err.Error())
		}
		slog.Info("Shutdown complete")
	}

	if cmd.Annotations["version"] == "testing" {
		doneChannel := make(chan struct{})
		go func() {
			slog.Info("Sleeping for 5 seconds")
			time.Sleep(5 * time.Second)
			slog.Info("Sending SIGTERM")
			stop(syscall.SIGTERM)
			doneChannel <- struct{}{}
		}()
		<-doneChannel
	} else {
		shutdown.AddWithParam(stop)
		shutdown.Listen(syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return nil
}
