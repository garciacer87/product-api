package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/garciacer87/product-api-challenge/internal/contract"
	"github.com/gorilla/mux"
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
func (s *server) getAll(w http.ResponseWriter, _ *http.Request) {
	prds, err := s.db.GetAll()
	if err != nil {
		logrus.Errorf("db error: %v", err)
		writeResponse(w, http.StatusInternalServerError, "could not get the list of products")
		return
	}

	if len(prds) > 0 {
		body, _ := json.Marshal(&prds)
		writeJSONResponse(w, http.StatusOK, body)
	} else {
		writeResponse(w, http.StatusNotFound, "No products found in database")
	}
}

func (s *server) get(w http.ResponseWriter, req *http.Request) {
	sku, ok := mux.Vars(req)["sku"]
	if !ok || sku == "" {
		logrus.Errorf("sku is not present")
		writeResponse(w, http.StatusBadRequest, "sku is not present")
		return
	}

	prd, err := s.db.Get(sku)
	if err != nil {
		logrus.Errorf("error retrieving product: %s", err)
		writeResponse(w, http.StatusInternalServerError, "could not retrieve product")
		return
	}

	if prd == nil {
		writeResponse(w, http.StatusNotFound, "product not found")
		return
	}

	body, _ := json.Marshal(&prd)
	writeJSONResponse(w, http.StatusOK, body)
}

func (s *server) update(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}

func (s *server) delete(w http.ResponseWriter, req *http.Request) {

	writeResponse(w, http.StatusOK, "ok!")
}
