package api

import (
	"fmt"

	"github.com/garciacer87/product-api-challenge/internal/contract"
)

type mockDB struct {
	throwError bool
	prdCount   int
}

func (mdb *mockDB) Create(prd contract.Product) error {
	if mdb.throwError {
		return fmt.Errorf("mocked error")
	}

	return nil
}

func (mdb *mockDB) GetAll() ([]contract.Product, error) {
	if mdb.throwError {
		return nil, fmt.Errorf("mock error")
	}

	prds := []contract.Product{}
	for i := 0; i < mdb.prdCount; i++ {
		prds = append(prds, getMockProduct())
	}

	return prds, nil
}

func (mdb *mockDB) Get(sku string) (*contract.Product, error) {
	if mdb.throwError && mdb.prdCount == 0 {
		return nil, fmt.Errorf("mocked error")
	}

	if mdb.prdCount == 0 {
		return nil, nil
	}

	prd := getMockProduct()

	return &prd, nil
}

func (mdb *mockDB) Update(prd contract.Product) error {
	if mdb.throwError {
		return fmt.Errorf("mocked error")
	}

	return nil
}

func (mdb *mockDB) Delete(sku string) error {
	if mdb.throwError {
		return fmt.Errorf("mocked error")
	}

	return nil
}

func (mdb *mockDB) Close() {}

func getMockProduct() contract.Product {
	return contract.Product{
		SKU:      "FAL-1000000",
		Name:     "name",
		Brand:    "brand",
		Size:     10,
		Price:    100.00,
		ImageURL: "http://aaaa",
		AltImages: []string{
			"http://bbbb",
			"http://cccc",
		},
	}
}
