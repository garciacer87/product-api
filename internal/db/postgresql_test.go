package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/garciacer87/product-api-challenge/internal/contract"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func initTestDB(t *testing.T) *migrate.Migrate {
	dbURI := "postgres://productapi:password@localhost:5432/productapitest"
	os.Setenv("DATABASE_URI", dbURI)

	m, err := migrate.New(
		"file://../../sql/postgresql",
		fmt.Sprintf("%s?sslmode=disable", dbURI),
	)

	if err != nil {
		t.Fatalf("could not create new migrate object: %s", err)
	}
	if err := m.Up(); err != nil {
		t.Fatalf("could not up migrate %s", err)
	}

	return m
}

func TestCreate(t *testing.T) {
	m := initTestDB(t)

	defer func() {
		if err := m.Down(); err != nil {
			t.Fatalf("could not down migrate %s", err)
		}
	}()

	db, err := NewPostgreSQLDB()
	if err != nil {
		t.Fatalf("could not init database connection: %s", err)
	}

	tests := map[string]struct {
		prd         contract.Product
		errExpected bool
	}{
		"#1: invalid sku": {
			prd:         contract.Product{SKU: "FAL-10000000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			errExpected: true,
		},
		"#2: invalid price": {
			prd:         contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100000000000000.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			errExpected: true,
		},
		"#3: valid case": {
			prd:         contract.Product{SKU: "FAL-1000000", Name: "name", Brand: "brand", Size: 10, Price: 100.00, ImageURL: "http://a", AltImages: []string{"http://b", "http://c"}},
			errExpected: false,
		},
	}

	for desc, tc := range tests {
		err = db.Create(tc.prd)
		isErr := err != nil
		if isErr != tc.errExpected {
			t.Fatalf("%s:\n got Error? %v.\n Error expected? %v.\n Error: %v", desc, isErr, tc.errExpected, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	m := initTestDB(t)
	defer func() {
		if err := m.Down(); err != nil {
			t.Fatalf("could not down migrate %s", err)
		}
	}()

	db, err := NewPostgreSQLDB()
	if err != nil {
		t.Fatalf("could not init database connection: %s", err)
	}

	prds, err := db.GetAll()
	if err != nil {
		t.Fatal()
	}

	if len(prds) != 0 {
		t.Errorf("#1: no products. Product slice must be empty")
	}

	prd := contract.Product{
		SKU: "FAL-1000000",
	}

	db.Create(prd)

	prds, err = db.GetAll()
	if err != nil {
		t.Fatal()
	}

	if len(prds) != 1 {
		t.Errorf("#2: Must be one product in the slice")
	}
}

func TestGet(t *testing.T) {
	m := initTestDB(t)
	defer func() {
		if err := m.Down(); err != nil {
			t.Fatalf("could not down migrate %s", err)
		}
	}()

	db, err := NewPostgreSQLDB()
	if err != nil {
		t.Fatalf("could not init database connection: %s", err)
	}

	db.Create(contract.Product{
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
	})

	tests := map[string]struct {
		sku      string
		prdFound bool
	}{
		"#1: valid case": {
			sku:      "FAL-1000000",
			prdFound: true,
		},
		"#2: product not found": {
			sku:      "FAL-123",
			prdFound: false,
		},
	}

	for desc, tc := range tests {
		prd, _ := db.Get(tc.sku)
		found := prd != nil

		if tc.prdFound != found {
			t.Errorf("%s:\n Product expected? %v\n Product found? %v", desc, tc.prdFound, found)
		}

		if prd != nil && prd.SKU != "FAL-1000000" {
			t.Errorf("%s\n expected different SKU", desc)
		}
	}

}
