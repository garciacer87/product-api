package api

import (
	"encoding/json"
	"net/http"

	"github.com/garciacer87/product-api/internal/contract"
)

//writeResponse writes reponse headers, code and body.
func writeResponse(w http.ResponseWriter, code int, msg interface{}) {
	resp := &contract.Response{
		Status:  code,
		Message: msg,
	}

	//convert the output to json
	respBody, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, code, respBody)
}

func writeJSONResponse(w http.ResponseWriter, code int, respBody []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBody)
}
