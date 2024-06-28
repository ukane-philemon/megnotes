package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ukane-philemon/megtask/webserver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	taskDB = "megTasks"
)

// Check that *MongoDB satisfies webserver.TaskDatabase.
var _ webserver.TaskDatabase = (*MongoDB)(nil)

// MongoDB implements webserver.TaskDatabase.
type MongoDB struct {
	db  *mongo.Database
	log *slog.Logger
}

// New connects to a mongo database and returns a new instance of *MongoDB.
func New(ctx context.Context, connectionURL string, logger *slog.Logger) (*MongoDB, error) {
	if connectionURL == "" {
		return nil, errors.New("missing mongodb database connection URL")
	}

	if logger == nil {
		return nil, errors.New("mongodb logger is required")
	}

	// Set server API version for the client.
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionURL).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect error: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("client.Ping error: %w", err)
	}

	logger.Info("Database has been connected and pinged successfully...")

	return &MongoDB{
		db:  client.Database(taskDB),
		log: logger,
	}, nil
}

// Shutdown attempts to shutdown the database.
func (mdb *MongoDB) Shutdown(ctx context.Context) error {
	client := mdb.db.Client()
	err := client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("client.Disconnect error: %w", err)
	}

	mdb.log.Info("Database has been shutdown successfully...")

	return nil
}
