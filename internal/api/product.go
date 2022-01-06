package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/garciacer87/product-api-challenge/internal/contract"
	"github.com/sirupsen/logrus"
)

var prodValidator *productValidator

func init() {
	prodValidator = newValidator()
}

//Handles the creation of the product. Decodes and validates the body. If everything is ok, inserts the new product into database and sends back a 200 reply
func (s *server) create(w http.ResponseWriter, req *http.Request) {
	prd := contract.Product{}

	err := json.NewDecoder(req.Body).Decode(&prd)
	if err != nil {
		logrus.Errorf("could not decode the body %v", err)
		writeResponse(w, http.StatusBadRequest, "could not decode the body")
		return
	}

	//validates product fields from decoded body
	err = prodValidator.Struct(prd)
	if err != nil {
		errs := prodValidator.translate(err)
		logrus.Printf("Validation error(s):\n%s", strings.Join(errs, " | "))
		writeResponse(w, http.StatusBadRequest, errs)
		return
	}

	//inserts the new product into database
	err = s.db.Create(prd)
	if err != nil {
		logrus.Errorf("db error: %s", err)
		writeResponse(w, http.StatusInternalServerError, "could not create new product")
		return
	}

	logrus.Infof("Product %s created", prd.SKU)
	writeResponse(w, http.StatusOK, "product successfully created")
}

//Handles the retrieval of all the products stored in the database
func (s *server) getAll(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}

func (s *server) get(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}

func (s *server) update(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}

func (s *server) delete(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}
