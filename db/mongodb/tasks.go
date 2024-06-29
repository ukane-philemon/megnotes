package mongodb

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/ukane-philemon/megtask/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTask creates a new task entry for a user.
func (mdb *MongoDB) CreateTask(userID string, taskDetail string) ([]*db.Task, error) {
	if userID == "" || taskDetail == "" {
		return nil, fmt.Errorf("%w: missing required argument(s)", db.ErrorInvalidRequest)
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

	return mdb.userTasks(userID, nil)
}

// Tasks returns all the tasks created by the provided userID.
func (mdb *MongoDB) Tasks(userID string) ([]*db.Task, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w: missing required argument", db.ErrorInvalidRequest)
	}

	return mdb.userTasks(userID, nil)
}

// TasksWithStatus returns user tasks that matches the provided filter.
func (mdb *MongoDB) TasksWithStatus(userID string, completed bool) ([]*db.Task, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w: missing required argument", db.ErrorInvalidRequest)
	}

	return mdb.userTasks(userID, bson.M{completedKey: completed})
}

// UpdateTask updates an existing task for the provided userID. If no task
// match the provided taskID, an ErrorInvalidRequest is returned.
func (mdb *MongoDB) UpdateTask(userID, taskID string, newTaskDetail string, markAsComplete *bool) ([]*db.Task, error) {
	nothingToUpdate := (newTaskDetail == "" && markAsComplete == nil)
	if userID == "" || taskID == "" || nothingToUpdate {
		return nil, fmt.Errorf("%w: missing required argument(s)", db.ErrorInvalidRequest)
	}

	taskDBID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("primitive.ObjectIDFromHex error: %w", err)
	}

	filter := bson.M{
		ownerIDKey: userID,
		dbIDKey:    taskDBID,
	}

	var task *dbTask
	err = mdb.tasksCollection.FindOne(mdb.ctx, filter).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: task does not exist", db.ErrorInvalidRequest)
		}
		return nil, fmt.Errorf("tasksCollection.FindOne error: %w", err)
	}

	if task.Completed {
		return nil, fmt.Errorf("%w: completed tasks cannot be updated", db.ErrorInvalidRequest)
	}

	update := make(bson.M, 0)
	if newTaskDetail != "" {
		update[taskDetailKey] = newTaskDetail
	}

	if markAsComplete != nil && *markAsComplete != false {
		update[completedKey] = *markAsComplete
	}

	res, err := mdb.tasksCollection.UpdateOne(mdb.ctx, filter, bson.M{"$set": update}, options.Update().SetUpsert(false))
	if err != nil {
		return nil, fmt.Errorf("tasksCollection.UpdateOne error: %w", err)
	}

	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("%w: task does not exist", db.ErrorInvalidRequest)
	}

	return mdb.userTasks(userID, nil)
}

// DeleteTask removes an existing task from the record of the user that
// match the provided userID. If no task match the provided taskID, an
// ErrorInvalidRequest is returned.
func (mdb *MongoDB) DeleteTask(userID, taskID string) ([]*db.Task, error) {
	if userID == "" || taskID == "" {
		return nil, fmt.Errorf("%w: missing required argument(s)", db.ErrorInvalidRequest)
	}

	taskDBID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("primitive.ObjectIDFromHex error: %w", err)
	}

	filter := bson.M{
		ownerIDKey: userID,
		dbIDKey:    taskDBID,
	}

	res, err := mdb.tasksCollection.DeleteOne(mdb.ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w: task does not exist", db.ErrorInvalidRequest)
		}
		return nil, fmt.Errorf("tasksCollection.FindOne error: %w", err)
	}

	if res.DeletedCount == 0 {
		return nil, fmt.Errorf("%w: task does not exist", db.ErrorInvalidRequest)
	}

	return mdb.userTasks(userID, nil)
}

// userTasks returns a list of tasks for the user with the provided userID.
// Tasks are sorted in descending order.
func (mdb *MongoDB) userTasks(userID string, extraFilter bson.M) ([]*db.Task, error) {
	filter := bson.M{ownerIDKey: userID}
	for key, val := range extraFilter {
		filter[key] = val
	}

	cur, err := mdb.tasksCollection.Find(mdb.ctx, filter)
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
