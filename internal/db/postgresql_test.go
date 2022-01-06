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
