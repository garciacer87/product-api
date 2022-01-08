package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/garciacer87/product-api/internal/contract"
)

func TestCreate(t *testing.T) {
	mockDB := &mockDB{}
	srv := NewServer("8081", mockDB)

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
		"#2: invalid sku": {
			prd:            contract.Product{SKU: "FAL-asdf", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#3: invalid sku": {
			prd:            contract.Product{SKU: "FAL-10", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#4: invalid name": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#5: blank name": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "   ", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#6: invalid price": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: -100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#7: invalid alternative images": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http/invalid-url"}},
			statusExpected: http.StatusBadRequest,
		},
		"#8: invalid url image": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http/invalid-url", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusBadRequest,
		},
		"#9: database error": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusInternalServerError,
		},
		"#10: valid case": {
			prd:            contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			statusExpected: http.StatusOK,
		},
	}

	for desc, tc := range tests {
		mockDB.throwError = tc.statusExpected == http.StatusInternalServerError

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
	tests := map[string]struct {
		srv            Server
		statusExpected int
		prdsExpected   int
	}{
		"#1: valid case":            {srv: NewServer("8081", &mockDB{prdCount: 2}), statusExpected: http.StatusOK, prdsExpected: 2},
		"#2: internal server error": {srv: NewServer("8081", &mockDB{throwError: true}), statusExpected: http.StatusInternalServerError},
		"#3: empty list":            {srv: NewServer("8081", &mockDB{}), statusExpected: http.StatusNotFound},
	}

	for desc, tc := range tests {
		go tc.srv.ListenAndServe()
		resp, err := http.Get("http://localhost:8081/product")
		if err != nil {
			t.Errorf("Error not expected: %v", err)
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n response code different than expected\n Got: %v\n Expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if resp.StatusCode == http.StatusOK {
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not get response body %v", err)
			}

			var prds []contract.Product
			err = json.Unmarshal(respBody, &prds)
			if err != nil {
				t.Fatalf("could not get response body %v", err)
			}

			if len(prds) != tc.prdsExpected {
				t.Errorf("there must be %v products in the slice", tc.prdsExpected)
			}

			resp.Body.Close()
		}

		if err := tc.srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}

func TestGet(t *testing.T) {
	tests := map[string]struct {
		db             *mockDB
		sku            string
		statusExpected int
	}{
		"#1: internal server error": {db: &mockDB{throwError: true, prdCount: 0}, sku: "FAL-1000001", statusExpected: http.StatusInternalServerError},
		"#2: product not found":     {db: &mockDB{prdCount: 0}, sku: "FAL-1000002", statusExpected: http.StatusNotFound},
		"#3: valid case":            {db: &mockDB{prdCount: 1}, sku: "FAL-1000003", statusExpected: http.StatusOK},
	}

	for desc, tc := range tests {
		srv := NewServer("8081", tc.db)
		go srv.ListenAndServe()

		url := fmt.Sprintf("http://localhost:8081/product/%s", tc.sku)
		resp, err := http.Get(url)

		if err != nil {
			t.Fatalf("error not expected")
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Status code got: %v\n Status code expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}

func TestUpdate(t *testing.T) {
	tests := map[string]struct {
		db             *mockDB
		sku            string
		patch          contract.Product
		statusExpected int
	}{
		"#1: internal server error": {db: &mockDB{throwError: true, prdCount: 1}, sku: "FAL-1000001", statusExpected: http.StatusInternalServerError},
		"#2: product not found":     {db: &mockDB{prdCount: 0}, sku: "FAL-1000002", statusExpected: http.StatusNotFound},
		"#3: valid case":            {db: &mockDB{prdCount: 1}, sku: "FAL-1000003", patch: contract.Product{Name: "new name"}, statusExpected: http.StatusOK},
	}

	for desc, tc := range tests {
		srv := NewServer("8081", tc.db)

		go srv.ListenAndServe()

		body, _ := json.Marshal(&tc.patch)
		url := fmt.Sprintf("http://localhost:8081/product/%s", tc.sku)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("error not expected")
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Status code got: %v\n Status code expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}

func TestDelete(t *testing.T) {
	tests := map[string]struct {
		db             *mockDB
		sku            string
		statusExpected int
	}{
		"#1: internal server error": {db: &mockDB{throwError: true, prdCount: 1}, sku: "FAL-1000001", statusExpected: http.StatusInternalServerError},
		"#2: product not found":     {db: &mockDB{prdCount: 0}, sku: "FAL-1000002", statusExpected: http.StatusNotFound},
		"#3: valid case":            {db: &mockDB{prdCount: 1}, sku: "FAL-1000003", statusExpected: http.StatusOK},
	}

	for desc, tc := range tests {
		srv := NewServer("8081", tc.db)
		go srv.ListenAndServe()

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

		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}
