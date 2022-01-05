package api

import (
	"encoding/json"
	"net/http"
)

func (s *server) healthHandler(w http.ResponseWriter, _ *http.Request) {
	resp, err := json.Marshal("healthy")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the content type to json for browsers
	w.Header().Set("Content-Type", "application/json")

	//Add the response code and response body.
	w.WriteHeader(http.StatusOK)

	w.Write(resp)
}
