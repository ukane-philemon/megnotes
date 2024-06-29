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
	s.writeJSONResponse(res, http.StatusOK, respBody)
}

// writeBadRequest writes an http.StatusBadRequest to the response header and an
// errorMessage.
func (s *WebServer) writeBadRequest(res http.ResponseWriter, errorMessage string) {
	s.writeJSONResponse(res, http.StatusBadRequest, map[string]string{
		"errorMessage": errorMessage,
	})
}

// writeServerError writes a server error and logs the provided error.
func (s *WebServer) writeServerError(res http.ResponseWriter, serverErr error) {
	s.log.Error("Server error: ", "err", serverErr)
	s.writeJSONResponse(res, http.StatusInternalServerError, map[string]string{
		"errorMessage": "Something unexpected happened, please try again later.",
	})
}

func (s *WebServer) writeJSONResponse(res http.ResponseWriter, code int, resp any) {
	res.WriteHeader(code)
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		s.log.Error("json.Marshal failed: ", "error", err)
	}
	res.Write(responseBytes)
}
