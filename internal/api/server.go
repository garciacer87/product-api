package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Server abstraction of a server
type Server interface {
	ListenAndServe() error
}

type server struct {
	httpPort   string
	httpServer *http.Server
}

//NewServer creates a new server object
func NewServer() Server {
	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "8080"
		logrus.Println("env variable PORT not defined. Using default port 8080")
	}

	srv := &server{
		httpPort: port,
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	product := r.PathPrefix("/product").Subrouter()
	product.HandleFunc("", create).Methods(http.MethodPost)
	product.HandleFunc("", getAll).Methods(http.MethodGet)
	product.HandleFunc("/{sku}", getAll).Methods(http.MethodGet)
	product.HandleFunc("/{sku}", update).Methods(http.MethodPatch)
	product.HandleFunc("/{sku}", delete).Methods(http.MethodDelete)

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

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeResponse(w, http.StatusOK, "healthy")
}
