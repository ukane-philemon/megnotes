package webserver

import (
	"context"

	"github.com/ukane-philemon/megtask/db"
)

type TaskDatabase interface {
	// CreateAccount creates a new user with the provided username and password.
	// An ErrorInvalidRequest will be returned is the username already exists.
	CreateAccount(username, password string) error
	// Login checks that the provided username and password matches a record in
	// the database and are correct. Returns ErrorInvalidRequest if the password
	// or username does not match any record.
	Login(username, password string) (*db.User, error)
	// CreateTask creates a new task entry for a user.
	CreateTask(userID string, taskDetail string) ([]*db.Task, error)
	// Tasks returns all the tasks created by the provided userID.
	Tasks(userID string) ([]*db.Task, error)
	// TasksWithStatus returns user tasks that matches the provided filter.
	TasksWithStatus(userID string, completed bool) ([]*db.Task, error)
	// UpdateTask updates an existing task for the provided userID. If no task
	// match the provided taskID, an ErrorInvalidRequest is returned.
	UpdateTask(userID, taskID string, newTaskDetail string) ([]*db.Task, error)
	// DeleteTask removes an existing task from the record of the user that
	// match the provided userID. If no task match the provided taskID, an
	// ErrorInvalidRequest is returned.
	DeleteTask(userID, taskID string) ([]*db.Task, error)
	// Shutdown gracefully disconnects the database after the server is
	// shutdown.
	Shutdown(ctx context.Context) error
}
