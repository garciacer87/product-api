package api

import (
	"encoding/json"
	"net/http"

	"github.com/garciacer87/product-api-challenge/internal/contract"
)

func create(w http.ResponseWriter, req *http.Request) {
	prd := &contract.Product{}

	err := json.NewDecoder(req.Body).Decode(&prd)
	if err != nil {
		writeErrorResponse(w, "could not decode the body")
	}

	writeResponse(w, http.StatusOK, "ok!")

}

func getAll(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}

func get(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}

func update(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}

func delete(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}
