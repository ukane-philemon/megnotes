package webserver

import (
	"fmt"
	"net/http"
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
