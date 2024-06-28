package mongodb

import "github.com/ukane-philemon/megtask/db"

// CreateTask creates a new task entry for a user.
func (mdb *MongoDB) CreateTask(userID string, taskDetail string) ([]*db.Task, error) {
	return nil, nil
}

// Tasks returns all the tasks created by the provided userID.
func (mdb *MongoDB) Tasks(userID string) ([]*db.Task, error) {
	return nil, nil
}

// UpdateTask updates an existing task for the provided userID. If no task
// match the provided taskID, an ErrorInvalidRequest is returned.
func (mdb *MongoDB) UpdateTask(userID, taskID string, newTaskDetail string) ([]*db.Task, error) {
	return nil, nil
}

// DeleteTask removes an existing task from the record of the user that
// match the provided userID. If no task match the provided taskID, an
// ErrorInvalidRequest is returned.
func (mdb *MongoDB) DeleteTask(userID, taskID string) ([]*db.Task, error) {
	return nil, nil
}

// TasksWithStatus returns user tasks that match the provided filter.
func (mdb *MongoDB) TasksWithStatus(userID string, completed bool) ([]*db.Task, error) {
	return nil, nil
}
