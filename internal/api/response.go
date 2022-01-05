package api

import (
	"encoding/json"
	"net/http"

	"github.com/garciacer87/product-api-challenge/internal/contract"
)

//writeBadRequestResponse writes an bad request response
func writeBadRequestResponse(w http.ResponseWriter, message string) {
	msg := contract.NewErrorResponse(http.StatusBadRequest, message)
	writeResponse(w, http.StatusBadRequest, msg)
}

//writeErrorResponse writes an error response
func writeErrorResponse(w http.ResponseWriter, message string) {
	msg := contract.NewErrorResponse(http.StatusInternalServerError, message)
	writeResponse(w, http.StatusInternalServerError, msg)
}

//writeNotFoundResponse writes a not found response
func writeNotFoundResponse(w http.ResponseWriter, message string) {
	msg := contract.NewErrorResponse(http.StatusNotFound, message)
	writeResponse(w, http.StatusNotFound, msg)
}

//writeResponse writes reponse headers, code and body.
func writeResponse(w http.ResponseWriter, code int, output interface{}) {
	//convert the output to json
	response, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the content type to json for browsers
	w.Header().Set("Content-Type", "application/json")

	//Add the response code and response body.
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		writeErrorResponse(w, "could not write the correct response")
	}
}
