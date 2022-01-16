package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SYSTEMTerror/GoHealth/cmd/app/middleware"
	"github.com/SYSTEMTerror/GoHealth/pkg/types"
	"github.com/gorilla/mux"
)

func (s *Server) handleSaveMedicine(w http.ResponseWriter, r *http.Request) {
	var item *types.Medicine
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := middleware.Authentication(r.Context())
	if err != nil {
		log.Println("app.medicines Save middleware.Authentication ERROR:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if id != 0 && !s.customersSvc.HasAnyRole(r.Context(), id, "ADMIN") {
		log.Println("app.managers handleManagerRegistration middleware.Authentication ERROR:", err)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	medicine, err := s.medicinesSvc.Save(r.Context(), item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsoner(w, medicine, http.StatusOK)
}

func (s *Server) handleGetMedicines(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	column, ok := vars["column"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	value, ok := vars["value"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	limitParam, ok := vars["limit"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	medicines, err := s.medicinesSvc.GetSomeMedicines(r.Context(), column, value, limit)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsoner(w, medicines, http.StatusOK)
}