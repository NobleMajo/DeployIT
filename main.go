package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/noblemajo/deployit/lib/bun"
)

var DisplayName string = "Unset"
var ShortName string = "unset"
var Version string = "?.?.?"
var Commit string = "???????"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	slog.Info("starting",
		"display_name", DisplayName,
		"version", Version,
		"commit", Commit,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := bun.RunEmbedded(ctx, slog.Default())
	if err != nil {
		slog.Error("bun runtime failed", "err", err)
		os.Exit(1)
	}

	slog.Info("done")
}
