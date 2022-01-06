package main

import (
	"github.com/garciacer87/product-api-challenge/internal/api"
	"github.com/garciacer87/product-api-challenge/internal/db"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := db.NewPostgreSQLDB()
	if err != nil {
		logrus.Panicf("could not initialize database: %v", err)
	}

	srv, err := api.NewServer(db)
	if err != nil {
		logrus.Panic(err)
	}

	srv.ListenAndServe()

	if err := srv.ListenAndServe(); err != nil {
		logrus.Panic(err)
	}
}
