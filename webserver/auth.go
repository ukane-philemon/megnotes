package webserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ukane-philemon/megtask/db"
)

// handleCreateAccount handles the "POST /create-account" endpoint and creates a
// new user account.
func (s *WebServer) handleCreateAccount(res http.ResponseWriter, req *http.Request) {
	form := new(usernameAndPassword)
	if !s.readPostBody(res, req, &form) {
		return
	}

	err := form.Validate()
	if err != nil {
		s.writeBadRequest(res, err.Error())
		return
	}

	err = s.taskDB.CreateAccount(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, db.ErrorInvalidRequest) {
			s.writeBadRequest(res, err.Error())
		} else {
			s.writeServerError(res, fmt.Errorf("taskDB.CreateAccount error: %w", err))
		}
		return
	}

	s.writeSuccess(res, map[string]string{
		"message": "Account created successfully, proceed to login.",
	})
}

// handleLogin handles the "POST /login" endpoint and attempts to logs a user
// into their account.
func (s *WebServer) handleLogin(res http.ResponseWriter, req *http.Request) {
	form := new(usernameAndPassword)
	if !s.readPostBody(res, req, &form) {
		return
	}

	err := form.Validate()
	if err != nil {
		s.writeBadRequest(res, err.Error())
		return
	}

	userInfo, err := s.taskDB.Login(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, db.ErrorInvalidRequest) {
			s.writeBadRequest(res, err.Error())
		} else {
			s.writeServerError(res, fmt.Errorf("taskDB.CreateAccount error: %w", err))
		}
		return
	}

	authToken, err := s.jwtManager.GenerateJWtToken(userInfo.ID)
	if err != nil {
		s.writeServerError(res, fmt.Errorf("jwtManager.GenerateJWtToken error: %w", err))
		return
	}

	s.writeSuccess(res, map[string]any{
		"userInfo":  userInfo,
		"authToken": authToken,
		"message":   "Login successful.",
	})
}
