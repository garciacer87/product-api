package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/garciacer87/product-api/docs"

	"github.com/garciacer87/product-api/internal/api"
	"github.com/garciacer87/product-api/internal/db"
	"github.com/sirupsen/logrus"
)

// @title Product-API
// @version 1.0.0
// @description Basic API to manage CRUD operations on products
// @contact.url https://github.com/garciacer87/product-api
// @host http://localhost:8080
// @BasePath /
func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "8080"
		logrus.Println("env variable PORT not defined. Using default port 8080")
	}

	dbURI := os.Getenv("DATABASE_URI")
	if dbURI == "" {
		logrus.Panicf("DATABASE_URI environment variable not defined")
	}

	db, err := db.NewPostgreSQLDB(dbURI)
	if err != nil {
		logrus.Panicf("could not initialize database: %v", err)
	}

	srv := api.NewServer(port, db)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()

	s := <-signalChan
	logrus.Infof("Signal triggered: %v", s)

	srv.Shutdown(context.Background())

	os.Exit(0)
}
