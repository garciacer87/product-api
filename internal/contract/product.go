package contract

type Product struct {
	SKU       string
	Name      string
	Brand     string
	Size      int
	Price     float64
	ImageURL  string
	AltImages []string
}
