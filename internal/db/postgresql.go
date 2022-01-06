package db

import (
	"context"
	"fmt"
	"os"

	"github.com/garciacer87/product-api-challenge/internal/contract"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgreSQLDB struct {
	pool *pgxpool.Pool
}

// NewPostgreSQLDB retrieves a new PostgreSQLDB object
func NewPostgreSQLDB() (Database, error) {
	databaseURL := os.Getenv("DATABASE_URI")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URI environment variable not defined")
	}

	pool, err := pgxpool.Connect(context.Background(), databaseURL)
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
func (p *PostgreSQLDB) Close() {
	logrus.Info("Closing database connections")
	p.pool.Close()
}

//Create inserts a new product
func (db *PostgreSQLDB) Create(prd contract.Product) error {
	query := `INSERT INTO public.product(sku, name, brand, size, price, image_url, alt_images) VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.pool.Exec(context.Background(), query, prd.SKU, prd.Name, prd.Brand, prd.Size, prd.Price, prd.ImageURL, prd.AltImages)
	if err != nil {
		return fmt.Errorf("could not create product: %v", err)
	}

	return nil
}

func (db *PostgreSQLDB) GetAll() ([]contract.Product, error) {
	return nil, nil
}

func (db *PostgreSQLDB) Get(sku string) (*contract.Product, error) {
	return nil, nil
}

func (db *PostgreSQLDB) Update(prd contract.Product) error {
	return nil
}

func (db *PostgreSQLDB) Delete(sku string) error {
	return nil
}
