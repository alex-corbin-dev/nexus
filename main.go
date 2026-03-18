package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alex-corbin-dev/nexus/internal/broker"
	"github.com/alex-corbin-dev/nexus/internal/config"
	"github.com/alex-corbin-dev/nexus/internal/telemetry"
)

var (
	configFile = flag.String("config", "nexus.yaml", "Path to configuration file")
	version    = flag.Bool("version", false, "Print version and exit")
)

const buildVersion = "0.4.2"

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("nexus %s\n", buildVersion)
		os.Exit(0)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load(*configFile)
	if err != nil {
		slog.Error("failed to load configuration", "error", err, "file", *configFile)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialise telemetry (metrics + tracing)
	shutdown, err := telemetry.Init(ctx, cfg.Observability)
	if err != nil {
		slog.Error("failed to initialise telemetry", "error", err)
		os.Exit(1)
	}
	defer shutdown(ctx)

	// Start the broker
	b, err := broker.New(cfg)
	if err != nil {
		slog.Error("failed to create broker", "error", err)
		os.Exit(1)
	}

	if err := b.Start(ctx); err != nil {
		slog.Error("failed to start broker", "error", err)
		os.Exit(1)
	}

	slog.Info("nexus broker started",
		"node_id", cfg.Cluster.NodeID,
		"listen_addr", cfg.Cluster.ListenAddr,
		"version", buildVersion,
	)

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	slog.Info("shutting down broker gracefully...")
	if err := b.Stop(ctx); err != nil {
		slog.Error("error during shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("broker stopped cleanly")
}
