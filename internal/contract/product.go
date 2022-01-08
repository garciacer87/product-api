package contract

//Product type used to represent a product entity
type Product struct {
	SKU       string   `json:"sku" validate:"required,sku"`
	Name      string   `json:"name" validate:"required,notblank,min=3,max=50"`
	Brand     string   `json:"brand" validate:"required,notblank,min=3,max=50"`
	Size      int      `json:"size" validate:"notblank,min=0,max=9999999999"`
	Price     float64  `json:"price" validate:"required,min=1.00,max=99999999.00"`
	ImageURL  string   `json:"imageURL" validate:"required,url"`
	AltImages []string `json:"altImages" validate:"altimages"`
}

//Patch set new values from patch to this product object
func (p *Product) Patch(patch Product) {
	if patch.SKU != "" {
		p.SKU = patch.SKU
	}

	if patch.Name != "" {
		p.Name = patch.Name
	}

	if patch.Brand != "" {
		p.Brand = patch.Brand
	}

	if patch.Size != 0 {
		p.Size = patch.Size
	}

	if patch.Price != 0.00 {
		p.Price = patch.Price
	}

	if patch.ImageURL != "" {
		p.ImageURL = patch.ImageURL
	}

	if len(patch.AltImages) > 0 {
		p.AltImages = patch.AltImages
	}
}
