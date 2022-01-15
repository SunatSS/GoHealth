package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/SYSTEMTerror/GoHealth/pkg/customers"
	"github.com/SYSTEMTerror/GoHealth/pkg/types"
)

type Response struct {
	Status     string `json:"status"`
	Reason     string `json:"reason,omitempty"`
	CustomerId int64  `json:"customerId,omitempty"`
}

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

	token, err := s.customersSvc.TokenForCustomer(r.Context(), item)
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

//handleValidateToken
func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	res := &Response{}

	var item *types.Token
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := s.customersSvc.AuthenticateCustomer(r.Context(), item)
	if errors.Is(err, customers.ErrNotFound) || errors.Is(err, customers.ErrInternal) {
		res.Status = "fail"
		res.Reason = "not found"
		
		jsoner(w, &res, http.StatusNotFound)
		return
	}
	if errors.Is(err, customers.ErrExpired) {
		res.Status = "fail"
		res.Reason = "expired"
		jsoner(w, &res, http.StatusBadRequest)
		return
	}

	res.Status = "ok"
	res.CustomerId = id
	jsoner(w, &res, http.StatusOK)
}