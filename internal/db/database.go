package db

import (
	"github.com/garciacer87/product-api-challenge/internal/contract"
)

//Database abstraction of database connection
type Database interface {
	Create(prd contract.Product) error
	GetAll() ([]contract.Product, error)
	Get(sku string) (*contract.Product, error)
	Update(prd contract.Product) error
	Delete(sku string) (*bool, error)
	Close()
}
