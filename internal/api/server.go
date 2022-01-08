package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/garciacer87/product-api/internal/db"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
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
func NewServer(port string, db db.Database) Server {
	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.PathPrefix("/swagger/{*}").Handler(httpSwagger.WrapHandler)

	srv := &server{
		httpPort: port,
		db:       db,
	}

	product := r.PathPrefix("/product").Subrouter()
	product.HandleFunc("", validateProduct(srv.create)).Methods(http.MethodPost)
	product.HandleFunc("", srv.getAll).Methods(http.MethodGet)
	product.HandleFunc("/{sku}", validateExistence(db, srv.get)).Methods(http.MethodGet)
	product.HandleFunc("/{sku}", validateExistence(db, validatePatchFields(db, srv.update))).Methods(http.MethodPatch)
	product.HandleFunc("/{sku}", validateExistence(db, srv.delete)).Methods(http.MethodDelete)

	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
		Handler: r,
	}

	return srv
}

//ListenAndServe starts the http server on the previously configurated port
func (s *server) ListenAndServe() error {
	logrus.Printf("serving on port %s\n", s.httpPort)
	return s.httpServer.ListenAndServe()
}

//Shutdown the http server
func (s *server) Shutdown(ctx context.Context) error {
	logrus.Infof("Shutting down API server")

	// close DB connection
	s.db.Close()

	return s.httpServer.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeResponse(w, http.StatusOK, "healthy")
}
