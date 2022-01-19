package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SYSTEMTerror/GoHealth/cmd/app/middleware"
	"github.com/SYSTEMTerror/GoHealth/pkg/customers"
	"github.com/SYSTEMTerror/GoHealth/pkg/medicines"
	"github.com/gorilla/mux"
)

//Server is structure for server with mux from gorilla/mux
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	medicinesSvc *medicines.Service
}

//NewServer creates new server with mux from gorilla/mux
func NewServer(mux *mux.Router, customersSvc *customers.Service, medicinesSvc *medicines.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc, medicinesSvc: medicinesSvc}
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
	customersSubrouter.HandleFunc("/edit", s.handleEditCustomer).Methods("POST")
	customersSubrouter.HandleFunc("/token", s.handleTokenForCustomer).Methods("POST")
	customersSubrouter.HandleFunc("/admin", s.handleMakeAdmin).Methods("POST")
	customersSubrouter.HandleFunc("/{id}", s.handleGetCustomerByID).Methods("GET")

	medicinesSubrouter := s.mux.PathPrefix("/api/medicines").Subrouter()
	medicinesSubrouter.Use(customersAuthenticateMd)

	medicinesSubrouter.HandleFunc("", s.handleSaveMedicine).Methods("POST")
	medicinesSubrouter.HandleFunc("/{column:(?:id|name|manafacturer|pharmacy_name)}/{value}/{limit}", s.handleGetMedicines).Methods("GET")
}

//function jsoner marshal interfaces to json and write to response writer
func jsoner(w http.ResponseWriter, v interface{}, code int) error {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("jsoner json.Marshal error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Println("jsoner w.Write error:", err)
		return err
	}
	return nil
}
