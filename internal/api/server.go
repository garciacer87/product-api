package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/garciacer87/product-api-challenge/internal/db"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Server abstraction of a server
type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type server struct {
	httpPort   string
	httpServer *http.Server
	db         db.Database
}

//NewServer creates a new server object
func NewServer(db db.Database) (Server, error) {
	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "8080"
		logrus.Println("env variable PORT not defined. Using default port 8080")
	}

	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	srv := &server{
		httpPort: port,
		db:       db,
	}

	product := r.PathPrefix("/product").Subrouter()
	product.HandleFunc("", srv.create).Methods(http.MethodPost)
	product.HandleFunc("", srv.getAll).Methods(http.MethodGet)
	product.HandleFunc("/{sku}", srv.get).Methods(http.MethodGet)
	product.HandleFunc("/{sku}", srv.update).Methods(http.MethodPatch)
	product.HandleFunc("/{sku}", srv.delete).Methods(http.MethodDelete)

	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
		Handler: r,
	}

	return srv, nil
}

//ListenAndServe starts the http server on the previously configurated port
func (s *server) ListenAndServe() error {
	logrus.Printf("serving on port %s\n", s.httpPort)
	return s.httpServer.ListenAndServe()
}

//Shutdown the http server
func (s *server) Shutdown(ctx context.Context) error {
	logrus.Infof("Shutting down API server")
	// shutdown server
	err := s.httpServer.Shutdown(ctx)

	// close DB connection
	s.db.Close()

	return err
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeResponse(w, http.StatusOK, "healthy")
}
