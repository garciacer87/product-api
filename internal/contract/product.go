package contract

//Product type used to represent a product entity
type Product struct {
	SKU       string   `json:"sku" validate:"required,sku"`
	Name      string   `json:"name" validate:"required,notblank,min=3,max=50"`
	Brand     string   `json:"brand" validate:"required,notblank,min=3,max=50"`
	Size      int      `json:"size" validate:"notblank"`
	Price     float64  `json:"price" validate:"required,min=1.00,max=99999999.00"`
	ImageURL  string   `json:"imageURL" validate:"required,url"`
	AltImages []string `json:"altImages" validate:"altimages"`
}
