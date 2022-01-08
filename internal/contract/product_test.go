package contract

import "testing"

func TestProductPatch(t *testing.T) {
	prd := Product{
		SKU:       "FAL-10000001",
		Name:      "old name",
		Brand:     "old brand",
		Size:      1,
		Price:     10.00,
		ImageURL:  "http://old",
		AltImages: []string{"http://old"},
	}

	patch := Product{
		SKU:       "FAL-10000001",
		Name:      "new name",
		Brand:     "new brand",
		Size:      2,
		Price:     20.00,
		ImageURL:  "http://new",
		AltImages: []string{"http://new"},
	}

	prd.Patch(patch)

	if prd.SKU != "FAL-10000001" {
		t.Errorf("sku different than expected: %v", prd.SKU)
	}

	if prd.Name != "new name" {
		t.Errorf("name different than expected: %v", prd.Name)
	}

	if prd.Brand != "new brand" {
		t.Errorf("brand different than expected: %v", prd.Brand)
	}

	if prd.Size != 2 {
		t.Errorf("size different than expected: %v", prd.Size)
	}

	if prd.Price != 20.00 {
		t.Errorf("price different than expected: %v", prd.Price)
	}

	if prd.ImageURL != "http://new" {
		t.Errorf("imageURL different than expected: %v", prd.ImageURL)
	}

	if prd.AltImages[0] != "http://new" {
		t.Errorf("altImages[0] different than expected: %v", prd.AltImages[0])
	}
}
