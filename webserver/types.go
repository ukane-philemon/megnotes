package webserver

import (
	"errors"
	"fmt"
	"regexp"
)

var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9]+$")

// usernameAndPassword is information required to create an account or login an
// existing user.
type usernameAndPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate ensures valid data is provided in createAccountRequest.
func (caq *usernameAndPassword) Validate() error {
	// Username cannot contain special characters.
	if caq.Username == "" || !usernameRegex.MatchString(caq.Username) {
		return errors.New("username can only contain alphanumeric characters")
	}

	const maxPassLength, minPassLength = 72, 6

	passLen := len(caq.Password)
	if passLen < minPassLength || passLen > maxPassLength {
		return fmt.Errorf("password must be more than %d characters but less than %d characters", minPassLength, maxPassLength)
	}

	return nil
}

// createTaskRequest is information required to create new task.
type createTaskRequest struct {
	TaskDetail string `json:"taskDetail"`
}

// updateTaskRequest is information that may be provided to update a task. Both cannot
// be empty.
type updateTaskRequest struct {
	TaskDetail      string `json:"taskDetail"` // optional
	MarkAsCompleted bool   `json:"markAsCompleted"`
}
