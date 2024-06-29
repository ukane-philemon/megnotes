package webserver

import (
	"context"
	"net/http"
)

const jwtHeader = "Megtask-Authentication-Token"
const userIDCtxKey = "userID"

// authMiddleware ensures the the correct and valid auth token is provided in
// this request.
func (s *WebServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get(jwtHeader)
		if authToken == "" {
			s.writeJSONResponse(res, http.StatusUnauthorized, "not authorized")
			return
		}

		userID, validToken := s.jwtManager.IsValidToken(authToken)
		if !validToken {
			s.writeJSONResponse(res, http.StatusUnauthorized, "not authorized")
			return
		}

		// Set userID for use in subsequent handlers.
		req = req.WithContext(context.WithValue(req.Context(), userIDCtxKey, userID))
		next.ServeHTTP(res, req)
	})
}

// reqUserID retrieves the userID from an authenticated request.
func (s *WebServer) reqUserID(req *http.Request) string {
	userIDVal := req.Context().Value(userIDCtxKey)
	if userIDVal != nil {
		return userIDVal.(string)
	}
	return ""
}
