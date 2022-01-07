package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/garciacer87/product-api-challenge/internal/contract"
)

const (
	errProdTag      = "error-product"
	notFoundProdtag = "product-not-found"
)

type mockDB struct{}

func (mdb *mockDB) Create(prd contract.Product) error {
	if prd.Name == errProdTag {
		return fmt.Errorf("mocked error!")
	}

	return nil
}

func (mdb *mockDB) GetAll() ([]contract.Product, error) {
	prds := []contract.Product{
		{SKU: "FAL-1000001"},
		{SKU: "FAL-1000002"},
	}

	return prds, nil
}

func (mdb *mockDB) Get(sku string) (*contract.Product, error) {
	switch sku {
	case errProdTag:
		return nil, fmt.Errorf("mocked error")
	case notFoundProdtag:
		return nil, nil
	default:
		return &contract.Product{SKU: sku}, nil
	}
}

func (mdb *mockDB) Update(prd contract.Product) error {
	return nil
}

func (mdb *mockDB) Delete(sku string) (*bool, error) {
	switch sku {
	case errProdTag:
		return nil, fmt.Errorf("mocked error")
	case notFoundProdtag:
		found := false
		return &found, nil
	default:
		found := true
		return &found, nil
	}
}

func (mdb *mockDB) Close() {}

func initTestServer(t *testing.T) Server {
	os.Setenv("PORT", "8081")

	db := &mockDB{}
	srv, err := NewServer(db)
	if err != nil {
		t.Fatalf("could not init the test server: %v", err)
	}

	return srv
}

func TestCreate(t *testing.T) {
	srv := initTestServer(t)

	defer func(srv Server) {
		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}(srv)

	go srv.ListenAndServe()

	tests := map[string]struct {
		prd            contract.Product
		statusExpected int
	}{
		"#1: invalid sku": {
			prd:            contract.Product{SKU: "-2345gsdfgsdgw345", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#2: invalid name": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#3: invalid price": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: -100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#4: invalid alternative images": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http/invalid-url"}},
			statusExpected: http.StatusBadRequest,
		},
		"#5: invalid url image": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http/invalid-url", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#6: database error": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: errProdTag, Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusInternalServerError,
		},
		"#7: valid case": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusOK,
		},
	}

	for desc, tc := range tests {
		body, _ := json.Marshal(&tc.prd)

		resp, err := http.Post("http://localhost:8081/product", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Errorf("Error not expected: %v", err)
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s Status expected: %v Status got: %v", desc, tc.statusExpected, resp.StatusCode)
		}
	}

}

func TestGetAll(t *testing.T) {
	srv := initTestServer(t)

	defer func(srv Server) {
		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}(srv)

	go srv.ListenAndServe()

	resp, err := http.Get("http://localhost:8081/product")
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code must be 200")
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not get response body %v", err)
	}

	var prds []contract.Product
	err = json.Unmarshal(respBody, &prds)
	if err != nil {
		t.Fatalf("could not get response body %v", err)
	}

	if len(prds) != 2 {
		t.Error("there must be 2 products in the slice")
	}
}

func TestGet(t *testing.T) {
	srv := initTestServer(t)

	defer func(srv Server) {
		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}(srv)

	go srv.ListenAndServe()

	tests := map[string]struct {
		sku            string
		statusExpected int
	}{
		"#1: internal server error": {sku: errProdTag, statusExpected: http.StatusInternalServerError},
		"#2: product not found":     {sku: notFoundProdtag, statusExpected: http.StatusNotFound},
		"#3: valid case":            {sku: "FAL-1000001", statusExpected: http.StatusOK},
	}

	for desc, tc := range tests {
		url := fmt.Sprintf("http://localhost:8081/product/%s", tc.sku)
		resp, err := http.Get(url)

		if err != nil {
			t.Fatalf("error not expected")
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Status code got: %v\n Status code expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}
	}
}

func TestDelete(t *testing.T) {
	srv := initTestServer(t)

	defer func(srv Server) {
		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}(srv)

	go srv.ListenAndServe()

	tests := map[string]struct {
		sku            string
		statusExpected int
	}{
		"#1: internal server error": {sku: errProdTag, statusExpected: http.StatusInternalServerError},
		"#2: product not found":     {sku: notFoundProdtag, statusExpected: http.StatusNotFound},
		"#3: valid case":            {sku: "FAL-1000001", statusExpected: http.StatusOK},
	}

	for desc, tc := range tests {
		url := fmt.Sprintf("http://localhost:8081/product/%s", tc.sku)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("error not expected")
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Status code got: %v\n Status code expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}
	}
}
