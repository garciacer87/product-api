package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/garciacer87/product-api/internal/contract"
	"github.com/garciacer87/product-api/internal/db"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//validates product fields
func validateProduct(next http.HandlerFunc) http.HandlerFunc {
	prodValidator := newValidator()

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logrus.Errorf("could not decode the body %v", err)
			writeResponse(w, http.StatusBadRequest, "could not decode the body")
			return
		}

		body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		prd := contract.Product{}

		err = json.NewDecoder(body).Decode(&prd)
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

		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		next(w, req)
	})
}

//validates product fields
func validateExistence(db db.Database, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		sku, ok := mux.Vars(req)["sku"]
		if !ok || sku == "" {
			logrus.Errorf("sku is not present")
			writeResponse(w, http.StatusBadRequest, "sku is not present")
			return
		}

		prd, err := db.Get(sku)
		if err != nil {
			logrus.Errorf("error retrieving product: %s", err)
			writeResponse(w, http.StatusInternalServerError, "could not retrieve product")
			return
		}

		if prd == nil {
			writeResponse(w, http.StatusNotFound, "product not found")
			return
		}

		next(w, req)
	})
}

//validates product fields from patch method
func validatePatchFields(db db.Database, next http.HandlerFunc) http.HandlerFunc {
	prodValidator := newValidator()

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logrus.Errorf("could not decode the body %v", err)
			writeResponse(w, http.StatusBadRequest, "could not decode the body")
			return
		}
		body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		sku := mux.Vars(req)["sku"]
		prd, _ := db.Get(sku)

		patch := contract.Product{}
		err = json.NewDecoder(body).Decode(&patch)
		if err != nil {
			logrus.Errorf("could not decode the body %v", err)
			writeResponse(w, http.StatusBadRequest, "could not decode the body")
			return
		}

		prd.Patch(patch)

		//validates product fields from decoded body
		err = prodValidator.Struct(prd)
		if err != nil {
			errs := prodValidator.translate(err)
			logrus.Printf("Validation error(s):\n%s", strings.Join(errs, " | "))
			writeResponse(w, http.StatusBadRequest, errs)
			return
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		next(w, req)
	})
}
