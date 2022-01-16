package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SYSTEMTerror/GoHealth/pkg/customers"
	"github.com/SYSTEMTerror/GoHealth/pkg/types"
)

//handleRegisterCustomer
func (s *Server) handleRegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var item *types.RegInfo
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	customer, err := s.customersSvc.RegisterCustomer(r.Context(), item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsoner(w, customer, http.StatusOK)
}

//handleTokenForCustomer
func (s *Server) handleTokenForCustomer(w http.ResponseWriter, r *http.Request) {
	var item *types.TokenInfo
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := s.customersSvc.Token(r.Context(), item)
	if err == customers.ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err == customers.ErrInvalidPassword {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsoner(w, token, http.StatusOK)
}

func (s *Server) handleEditCustomer(w http.ResponseWriter, r *http.Request) {
	var item *types.Customer
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	customer, err := s.customersSvc.EditCustomer(r.Context(), item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsoner(w, customer, http.StatusOK)
}