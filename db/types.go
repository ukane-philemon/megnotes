package db

import (
	"errors"
)

var (
	ErrorInvalidRequest = errors.New("invalid request")
)

// User is information about a user.
type User struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Tasks    []*Task `json:"tasks"`
}

// Task is information about a user's task item.
type Task struct {
	ID string `json:"id"`
	TaskInfo
}

type TaskInfo struct {
	Detail    string `json:"detail"`
	Completed bool   `json:"completed"`
	Timestamp int64  `json:"timestamp"`
}
