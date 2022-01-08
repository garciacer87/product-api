package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestValidatePatch(t *testing.T) {
	tests := map[string]struct {
		patch          []byte
		sku            string
		statusExpected int
	}{
		"1: invalid sku":  {patch: []byte(`{"name":"a"}`), sku: "FAL-10000001", statusExpected: http.StatusBadRequest},
		"2: invalid body": {patch: []byte(`{//"name":"name"}`), sku: "FAL-10000001", statusExpected: http.StatusBadRequest},
	}

	for desc, tc := range tests {
		srv := NewServer("8081", &mockDB{prdCount: 1})

		go srv.ListenAndServe()

		url := fmt.Sprintf("http://localhost:8081/product/%s", tc.sku)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(tc.patch))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("error not expected")
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Got: %v\n Expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}

func TestValidateProduct(t *testing.T) {
	tests := map[string]struct {
		prd            []byte
		statusExpected int
	}{
		"1: invalid sku":  {prd: []byte(`{"sku":"FAL-10000001", "name":"a"}, "brand":"brand", "size":10`), statusExpected: http.StatusBadRequest},
		"2: invalid body": {prd: []byte(`{//"name":"name"}`), statusExpected: http.StatusBadRequest},
	}

	for desc, tc := range tests {
		srv := NewServer("8081", &mockDB{prdCount: 1})

		go srv.ListenAndServe()

		resp, err := http.Post("http://localhost:8081/product", "application/json", bytes.NewBuffer(tc.prd))
		if err != nil {
			t.Errorf("Error not expected: %v", err)
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Got: %v\n Expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}
