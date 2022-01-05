package main

import (
	"github.com/garciacer87/product-api-challenge/internal/api"
	"github.com/sirupsen/logrus"
)

func main() {
	srv := api.NewServer()
	srv.ListenAndServe()

	if err := srv.ListenAndServe(); err != nil {
		logrus.Println(err)
	}
}
