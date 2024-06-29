package mongodb

import (
	"fmt"
	"sort"

	"github.com/ukane-philemon/megtask/db"
	"go.mongodb.org/mongo-driver/bson"
)

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

// userTasks returns a list of tasks for the user with the provided userID.
// Tasks are sorted in descending order.
func (mdb *MongoDB) userTasks(userID string) ([]*db.Task, error) {
	cur, err := mdb.tasksCollection.Find(mdb.ctx, bson.M{ownerIDKey: userID})
	if err != nil {
		return nil, fmt.Errorf("tasksCollection.Find error: %w", err)
	}

	var dbTasks []*dbTask
	err = cur.Decode(&dbTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to decode retrieved tasks: %w", err)
	}

	userTasks := make([]*db.Task, 0, len(dbTasks))
	for _, task := range dbTasks {
		userTasks = append(userTasks, &db.Task{
			ID: task.ID.Hex(),
			TaskInfo: db.TaskInfo{
				Detail:    task.Detail,
				Completed: task.Completed,
				Timestamp: task.Timestamp,
			},
		})
	}

	sort.SliceStable(userTasks, func(i, j int) bool {
		return userTasks[i].Timestamp > userTasks[j].Timestamp
	})

	return userTasks, nil
}
