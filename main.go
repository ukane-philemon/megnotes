package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ukane-philemon/megtask/db/mongodb"
	"github.com/ukane-philemon/megtask/webserver"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var dbConnectionURL string
	flag.StringVar(&dbConnectionURL, "dbURL", "", "dbConnectionURL is a mongoDB connection URL and must be provided to connect to a database.")
	flag.Parse()

	logger := slog.New(slog.Default().Handler())

	// Connect to database.
	db, err := mongodb.New(ctx, dbConnectionURL, logger)
	if err != nil {
		println("mongodb.New error: ", err.Error())
		os.Exit(1)
	}

	// Ensure graceful shutdown by capturing SIGINT and SIGTERM signals.
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdownChan
		// cancel the context so server can shutdown gracefully.
		cancel()
	}()

	server := webserver.New(db, logger)
	err = server.Start(ctx)
	if err != nil {
		println("webServer exited with unexpected error: ", err.Error())
		return
	}
}
