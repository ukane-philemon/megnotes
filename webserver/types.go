package webserver

import (
	"errors"
	"fmt"
	"regexp"
)

var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9]+$")

// createAccountRequest is the information required to create and account.
type createAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate ensures valid data is provided in createAccountRequest.
func (caq *createAccountRequest) Validate() error {
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
