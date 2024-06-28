package webserver

import (
	"encoding/json"
	"net/http"
)

// readPostBody reads the request body into body.
func (s *WebServer) readPostBody(res http.ResponseWriter, req *http.Request, body any) bool {
	err := json.NewDecoder(req.Body).Decode(body)
	if err != nil {
		s.writeBadRequest(res, "Invalid request body")
		return false
	}
	return true
}

// writeSuccess writes a success response.
func (s *WebServer) writeSuccess(res http.ResponseWriter, respBody any) {
	res.WriteHeader(http.StatusOK)
	responseBytes, err := json.Marshal(respBody)
	if err != nil {
		s.log.Error("json.Marshal failed to send response: ", "error", err)
	}
	res.Write(responseBytes)
}

// writeBadRequest writes an http.StatusBadRequest to the response header and an
// errorMessage.
func (s *WebServer) writeBadRequest(res http.ResponseWriter, errorMessage string) {
	res.WriteHeader(http.StatusBadRequest)
	responseBytes, err := json.Marshal(map[string]string{
		"errorMessage": errorMessage,
	})
	if err != nil {
		s.log.Error("json.Marshal failed to send bad request: ", "error", err)
	}
	res.Write(responseBytes)
}

// writeServerError writes a server error and logs the provided error.
func (s *WebServer) writeServerError(res http.ResponseWriter, serverErr error) {
	s.log.Error("Server error: ", "err", serverErr)

	res.WriteHeader(http.StatusInternalServerError)
	responseBytes, err := json.Marshal(map[string]string{
		"errorMessage": "Something unexpected happened. Sorry, try again.",
	})
	if err != nil {
		s.log.Error("json.Marshal failed to send server error: ", "error", err)
	}
	res.Write(responseBytes)
}
