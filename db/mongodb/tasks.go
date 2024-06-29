package mongodb

import (
	"fmt"
	"sort"
	"time"

	"github.com/ukane-philemon/megtask/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTask creates a new task entry for a user.
func (mdb *MongoDB) CreateTask(userID string, taskDetail string) ([]*db.Task, error) {
	if userID == "" || taskDetail == "" {
		return nil, fmt.Errorf("%w: missing required arguments", db.ErrorInvalidRequest)
	}

	userDBID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("primitive.ObjectIDFromHex error: %w", err)
	}

	// Check if user really exists.
	filter := bson.M{dbIDKey: userDBID}
	nUsersFound, err := mdb.usersCollection.CountDocuments(mdb.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("usersCollection.CountDocuments error: %w", err)
	}

	if nUsersFound != 1 {
		return nil, fmt.Errorf("expected userID to match one user, got %d", nUsersFound)
	}

	taskInfo := &dbTask{
		ID:      primitive.NewObjectID(),
		OwnerID: userID,
		TaskInfo: db.TaskInfo{
			Detail:    taskDetail,
			Timestamp: time.Now().Unix(),
		},
	}

	_, err = mdb.tasksCollection.InsertOne(mdb.ctx, taskInfo)
	if err != nil {
		return nil, fmt.Errorf("tasksCollection.InsertOne error: %w", err)
	}

	return mdb.userTasks(userID)
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
	err = cur.All(mdb.ctx, &dbTasks)
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
