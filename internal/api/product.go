package api

import (
	"encoding/json"
	"net/http"

	"github.com/garciacer87/product-api-challenge/internal/contract"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Handles the creation of the product. Decodes and validates the body. If everything is ok, inserts the new product into database and sends back a 200 reply
func (s *server) create(w http.ResponseWriter, req *http.Request) {
	prd := contract.Product{}
	json.NewDecoder(req.Body).Decode(&prd)

	//inserts the new product into database
	err := s.db.Create(prd)
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

//Handles the GET method in order to get a product by sku
func (s *server) get(w http.ResponseWriter, req *http.Request) {
	sku := mux.Vars(req)["sku"]

	prd, _ := s.db.Get(sku)

	body, _ := json.Marshal(&prd)

	writeJSONResponse(w, http.StatusOK, body)
}

//Handles the PATCH method in order to update a product
func (s *server) update(w http.ResponseWriter, req *http.Request) {
	var (
		sku   string           = mux.Vars(req)["sku"]
		patch contract.Product = contract.Product{}
	)

	prd, _ := s.db.Get(sku)

	json.NewDecoder(req.Body).Decode(&patch)

	prd.Patch(patch)

	err := s.db.Update(*prd)
	if err != nil {
		logrus.Errorf("error updating product: %s", err)
		writeResponse(w, http.StatusInternalServerError, "could not update product")
		return
	}

	writeResponse(w, http.StatusOK, "product successfully updated")
}

//Handles the DELETE method in order to delete a product
func (s *server) delete(w http.ResponseWriter, req *http.Request) {
	sku := mux.Vars(req)["sku"]

	err := s.db.Delete(sku)
	if err != nil {
		logrus.Errorf("error deleting product: %s", err)
		writeResponse(w, http.StatusInternalServerError, "could not delete product")
		return
	}

	writeResponse(w, http.StatusOK, "product successfully deleted")
}
