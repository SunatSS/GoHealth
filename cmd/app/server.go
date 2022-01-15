package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SYSTEMTerror/GoHealth/cmd/app/middleware"
	"github.com/SYSTEMTerror/GoHealth/pkg/customers"
	"github.com/gorilla/mux"
)

//Server is structure for server with mux from net/http
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
}

//NewServer creates new server with mux from net/http
func NewServer(mux *mux.Router, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc}
}

// ServeHTTP
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init initializes server
func (s *Server) Init() {
	s.mux.Use(middleware.Logger)

	customersAuthenticateMd := middleware.Authenticate(s.customersSvc.IDByToken)

	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customersAuthenticateMd)

	customersSubrouter.HandleFunc("", s.handleRegisterCustomer).Methods("POST")
	customersSubrouter.HandleFunc("/token", s.handleTokenForCustomer).Methods("POST")
	customersSubrouter.HandleFunc("/token/validate", s.handleValidateToken).Methods("POST")
}

//function jsoner marshal interface to json and write to response writer
func jsoner(w http.ResponseWriter, v interface{}, code int) error {
	data, err := json.Marshal(v)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}