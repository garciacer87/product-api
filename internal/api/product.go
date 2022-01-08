package api

import (
	"encoding/json"
	"net/http"

	"github.com/garciacer87/product-api/internal/contract"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// create godoc
// @Summary Creates a new product
// @Description Creates a new product
// @Tags product create
// @Accept json
// @Success 200 {object} contract.Response{status=int,message=object}
// @Failure 400,500 {object} contract.Response{status=int,message=object}
// @Param product body contract.Product true "product"
// @Router /product [post]
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

// getAll godoc
// @Summary Retrieves all the products stored in the database
// @Description Retrieves all the products stored in the database
// @Tags product list
// @Success 200 {array} contract.Product
// @Failure 404,500 {object} contract.Response{status=int,message=object}
// @Router /product [get]
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

// get godoc
// @Summary Get a product by its SKU
// @Description Get a product by its SKU
// @Tags product get
// @Accept json
// @Success 200 {object} contract.Product
// @Failure 400,404,500 {object} contract.Response{status=int,message=object}
// @Param sku path string true "product sku"
// @Router /product/{sku} [get]
func (s *server) get(w http.ResponseWriter, req *http.Request) {
	sku := mux.Vars(req)["sku"]

	prd, _ := s.db.Get(sku)

	body, _ := json.Marshal(&prd)

	writeJSONResponse(w, http.StatusOK, body)
}

// update godoc
// @Summary Updates an existing product
// @Description Updates a existing product
// @Tags product patch
// @Accept json
// @Success 200 {object} contract.Response{status=int,message=object}
// @Failure 400,404,500 {object} contract.Response{status=int,message=object}
// @Param sku path string true "product sku"
// @Param patch body contract.Product true "product patch"
// @Router /product/{sku} [patch]
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

// delete godoc
// @Summary Deletes an existing product
// @Description Deletes a existing product
// @Tags product delete
// @Success 200 {object} contract.Response{status=int,message=object}
// @Failure 400,404,500 {object} contract.Response{status=int,message=object}
// @Param sku path string true "sku product"
// @Param patch body contract.Product true "product patch"
// @Router /product/{sku} [delete]
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
