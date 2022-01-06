package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/garciacer87/product-api-challenge/internal/contract"
)

type mockDB struct{}

func (mdb *mockDB) Create(prd contract.Product) error {
	return nil
}

func (mdb *mockDB) GetAll() ([]contract.Product, error) {
	return nil, nil
}

func (mdb *mockDB) Get(sku string) (*contract.Product, error) {
	return nil, nil
}

func (mdb *mockDB) Update(prd contract.Product) error {
	return nil
}

func (mdb *mockDB) Delete(sku string) error {
	return nil
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
		"#6: valid case": {
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
