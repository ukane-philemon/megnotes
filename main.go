package main

import (
	"context"
	"log/slog"

	"github.com/ukane-philemon/megnotes/webserver"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.Default().Handler())
	server := webserver.New(logger)
	err := server.Start(ctx)
	if err != nil {
		logger.Error("webServer exited with unexpected error", "msg", err)
	}
}
