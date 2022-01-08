package db

import (
	"context"
	"fmt"

	"github.com/garciacer87/product-api-challenge/internal/contract"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

//PostgreSQLDB implementation of postgresql database
type PostgreSQLDB struct {
	pool *pgxpool.Pool
}

// NewPostgreSQLDB retrieves a new PostgreSQLDB object
func NewPostgreSQLDB(dbURI string) (Database, error) {
	pool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		return nil, fmt.Errorf("could not create database connection: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database is unreachable: %v", err)
	}
	logrus.Info("Connection succesfully to database")

	return &PostgreSQLDB{pool}, nil
}

// Close close connections from the pool
func (db *PostgreSQLDB) Close() {
	logrus.Info("Closing database connections")
	db.pool.Close()
}

//Create inserts a new product
func (db *PostgreSQLDB) Create(prd contract.Product) error {
	query := "INSERT INTO public.product(sku, name, brand, size, price, image_url, alt_images) VALUES($1, $2, $3, $4, $5, $6, $7)"

	_, err := db.pool.Exec(context.Background(), query, prd.SKU, prd.Name, prd.Brand, prd.Size, prd.Price, prd.ImageURL, prd.AltImages)
	if err != nil {
		return fmt.Errorf("could not create product: %v", err)
	}

	return nil
}

//GetAll retrieves a slice of the products stored in database
func (db *PostgreSQLDB) GetAll() ([]contract.Product, error) {
	query := "SELECT sku, name, brand, size, price, image_url, alt_images FROM public.product"

	var (
		sku, name, brand, imageURL string
		size                       int
		price                      float64
		altImages                  []string
	)

	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("could not get products: %v", err)
	}
	defer rows.Close()

	prds := make([]contract.Product, 0)
	for rows.Next() {
		if err = rows.Scan(&sku, &name, &brand, &size, &price, &imageURL, &altImages); err != nil {
			return nil, fmt.Errorf("could not get products: %v", err)
		}
		prds = append(prds, contract.Product{
			SKU:       sku,
			Name:      name,
			Brand:     brand,
			Size:      size,
			Price:     price,
			ImageURL:  imageURL,
			AltImages: altImages,
		})
	}

	return prds, nil
}

//Get retrieves a product by its SKU
func (db *PostgreSQLDB) Get(sku string) (*contract.Product, error) {
	query := "SELECT name, brand, size, price, image_url, alt_images FROM public.product WHERE sku = $1"

	var (
		name, brand, imageURL string
		size                  int
		price                 float64
		altImages             []string
	)

	row := db.pool.QueryRow(context.Background(), query, sku)
	err := row.Scan(&name, &brand, &size, &price, &imageURL, &altImages)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, nil
		default:
			return nil, fmt.Errorf("could not get product: %v", err)
		}
	}

	return &contract.Product{
		SKU:       sku,
		Name:      name,
		Brand:     brand,
		Size:      size,
		Price:     price,
		ImageURL:  imageURL,
		AltImages: altImages,
	}, nil
}

//Update updates a product by its SKU
func (db *PostgreSQLDB) Update(prd contract.Product) error {
	query := "UPDATE public.product SET name=$1, brand=$2, size=$3, price=$4, image_url=$5, alt_images=$6 WHERE sku=$7"

	_, err := db.pool.Exec(context.Background(), query, prd.Name, prd.Brand, prd.Size, prd.Price, prd.ImageURL, prd.AltImages, prd.SKU)
	if err != nil {
		return fmt.Errorf("could not update product: %v", err)
	}

	return nil
}

//Delete deletes product by its SKU
func (db *PostgreSQLDB) Delete(sku string) error {
	query := `DELETE FROM public.product WHERE sku=$1`
	_, err := db.pool.Exec(context.Background(), query, sku)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}

	return nil
}
