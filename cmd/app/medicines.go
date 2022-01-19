package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SYSTEMTerror/GoHealth/cmd/app/middleware"
	"github.com/SYSTEMTerror/GoHealth/pkg/customers"
	"github.com/SYSTEMTerror/GoHealth/pkg/types"
	"github.com/gorilla/mux"
)

//handleSaveMedicine
func (s *Server) handleSaveMedicine(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.Authentication(r.Context())
	if err != nil {
		log.Println("handleSaveMedicine middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	isAdmin, err := s.customersSvc.IsAdmin(r.Context(), id)
	if err == customers.ErrNotFound {
		log.Println("handleSaveMedicine s.customersSvc.IsAdmin Not Found:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("handleSaveMedicine s.customersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !isAdmin {
		log.Println("handleSaveMedicine s.customersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var item *types.Medicine
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println("handleSaveMedicine json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	medicine, err := s.medicinesSvc.Save(r.Context(), item)
	if err != nil {
		log.Println("handleSaveMedicine s.medicinesSvc.Save error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsoner(w, medicine, http.StatusOK)
}

//handleGetMedicines
func (s *Server) handleGetMedicines(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	column, ok := vars["column"]
	if !ok {
		log.Println("handleGetMedicines mux.Vars column:", ok)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	value, ok := vars["value"]
	if !ok {
		log.Println("handleGetMedicines mux.Vars value:", ok)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	limitParam, ok := vars["limit"]
	if !ok {
		log.Println("handleGetMedicines mux.Vars limit:", ok)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		log.Println("handleGetMedicines strconv.Atoi error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	medicines, err := s.medicinesSvc.GetSomeMedicines(r.Context(), column, value, limit)
	if err == customers.ErrNotFound {
		log.Println("handleGetMedicines s.medicinesSvc.GetSomeMedicines Not Found:", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("handleGetMedicines s.medicinesSvc.GetSomeMedicines error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jsoner(w, medicines, http.StatusOK)
}
