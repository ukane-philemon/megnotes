package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ukane-philemon/megtask/db"
)

const (
	// taskStatusQueryKey is the expected query key to provide a task filter.
	taskStatusQueryKey = "status"

	// completedTasksFilter is the status filter used to filter only completed
	// tasks.
	completedTasksFilter = "completed"
	// pendingTasksFilter is the status filter used to filter only pending
	// tasks.
	pendingTasksFilter = "pending"
)

// handleCreateTask handles the "POST /task" endpoint and creates a new task
// entry for a user.
func (s *WebServer) handleCreateTask(res http.ResponseWriter, req *http.Request) {
	form := new(createTaskRequest)
	if !s.readPostBody(res, req, &form) {
		return
	}

	if form.TaskDetail == "" {
		s.writeBadRequest(res, "missing task detail")
		return
	}

	userID := s.reqUserID(req)
	userTasks, err := s.taskDB.CreateTask(userID, form.TaskDetail)
	if err != nil {
		s.writeServerError(res, fmt.Errorf("taskDB.CreateTask error: %w", err))
		return
	}

	s.writeSuccess(res, map[string]any{
		"tasks": userTasks,
	})
}

// handleRetrieveTasks handles the "GET /tasks" endpoint and returns all users
// tasks sorted by timestamp. This endpoint excepts an optional "status" query
// parameter that can either be "pending" or "completed".
func (s *WebServer) handleRetrieveTasks(res http.ResponseWriter, req *http.Request) {
	status := req.URL.Query().Get(taskStatusQueryKey)
	if status != "" && !strings.EqualFold(status, pendingTasksFilter) && !strings.EqualFold(status, completedTasksFilter) {
		s.writeBadRequest(res, `"status" query param can either be "pending" or "completed"`)
		return
	}

	userID := s.reqUserID(req)

	var userTasks []*db.Task
	var err error
	var methodName string
	if status != "" {
		methodName = "taskDB.TasksWithStatus"
		filterCompleted := strings.EqualFold(status, completedTasksFilter)
		userTasks, err = s.taskDB.TasksWithStatus(userID, filterCompleted)
	} else {
		methodName = "taskDB.Tasks"
		userTasks, err = s.taskDB.Tasks(userID)
	}
	if err != nil {
		s.writeServerError(res, fmt.Errorf("%s error: %w", methodName, err))
		return
	}

	s.writeSuccess(res, map[string]any{
		"tasks": userTasks,
	})
}

// handleUpdateTask handles the "PATCH /task" endpoint and updates an existing
// task. A task that has been completed cannot be updated back to pending.
func (s *WebServer) handleUpdateTask(res http.ResponseWriter, req *http.Request) {
	form := new(updateTaskRequest)
	if !s.readPostBody(res, req, &form) {
		return
	}

	if form.TaskID == "" || (form.TaskDetail == "" && !form.MarkAsCompleted) {
		s.writeBadRequest(res, "missing required data")
		return
	}

	var markAsCompleted *bool
	if form.MarkAsCompleted {
		markAsCompleted = &form.MarkAsCompleted
	}

	userID := s.reqUserID(req)

	userTasks, err := s.taskDB.UpdateTask(userID, form.TaskID, form.TaskDetail, markAsCompleted)
	if err != nil {
		if errors.Is(err, db.ErrorInvalidRequest) {
			s.writeBadRequest(res, err.Error())
		} else {
			s.writeServerError(res, fmt.Errorf("taskDB.UpdateTask: %w", err))
		}
		return
	}

	s.writeSuccess(res, map[string]any{
		"tasks": userTasks,
	})
}
