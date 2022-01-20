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

//handleRegisterCustomer
func (s *Server) handleRegisterCustomer(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleRegisterCustomer started")

	var item *types.RegInfo
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		loggers.ErrorLogger.Println("handleRegisterCustomer json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	customer, statusCode, err := s.customersSvc.RegisterCustomer(r.Context(), item)
	if err != nil {
		loggers.ErrorLogger.Println("handleRegisterCustomer s.customersSvc.RegisterCustomer error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, customer, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleRegisterCustomer jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleRegisterCustomer finished with any error!")
}

//handleTokenForCustomer
func (s *Server) handleTokenForCustomer(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleTokenForCustomer started")

	var item *types.TokenInfo
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		loggers.ErrorLogger.Println("handleTokenForCustomer json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, statusCode, err := s.customersSvc.Token(r.Context(), item)
	if err != nil {
		loggers.ErrorLogger.Println("handleTokenForCustomer s.customersSvc.Token error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, token, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleTokenForCustomer jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleTokenForCustomer finished with any error!")
}

//handleEditCustomer
func (s *Server) handleEditCustomer(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleEditCustomer started")

	var item *types.Customer
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		loggers.ErrorLogger.Println("handleTokenForCustomer json.NewDecoder error:")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleTokenForCustomer middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	item.ID = id

	statusCode, err := s.customersSvc.EditCustomer(r.Context(), item)
	if err != nil {
		loggers.ErrorLogger.Println("handleTokenForCustomer s.customersSvc.Token error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, item, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleTokenForCustomer jsoner error:", err)
		return
	}
	loggers.InfoLogger.Println("handleEditCustomer finished with any error!")
}

//handleMakeAdmin makes a customer with id an admin
func (s *Server) handleMakeAdmin(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleMakeAdmin started")

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	isAdmin, statusCode, err := s.customersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.customersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	if !isAdmin {
		loggers.ErrorLogger.Println("handleMakeAdmin s.customersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var makeAdminInfo *types.MakeAdminInfo
	err = json.NewDecoder(r.Body).Decode(&makeAdminInfo)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin json.NewDecoder error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	statusCode, err = s.customersSvc.MakeAdmin(r.Context(), makeAdminInfo)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.customersSvc.MakeAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, makeAdminInfo, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin jsoner error:", err)
		return
	}
}

//handleGetCustomerByID
func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleGetCustomerByID started")

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	isAdmin, statusCode, err := s.customersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.customersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	if !isAdmin {
		loggers.ErrorLogger.Println("handleMakeAdmin s.customersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		loggers.ErrorLogger.Println("handleGetCustomerByID mux.Vars(r) ID not found")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetCustomerByID strconv.ParseInt error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	customer, statusCode, err := s.customersSvc.GetCustomerByID(r.Context(), id)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetCustomerByID s.customersSvc.GetCustomerByID error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	err = jsoner(w, customer, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetCustomerByID jsoner error:", err)
		return
	}
}

//handleGetAllCustomers
func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	loggers, err := middleware.GetLoggers(r.Context())
	if err != nil {
		log.Println("LOGGERS DON'T WORK!!!")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loggers.InfoLogger.Println("handleGetAllCustomers started")

	adminId, err := middleware.Authentication(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin middleware.Authentication error:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	isAdmin, statusCode, err := s.customersSvc.IsAdmin(r.Context(), adminId)
	if err != nil {
		loggers.ErrorLogger.Println("handleMakeAdmin s.customersSvc.IsAdmin error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	if !isAdmin {
		loggers.ErrorLogger.Println("handleMakeAdmin s.customersSvc.IsAdmin isAdmin:", isAdmin)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	customersArr, statusCode, err := s.customersSvc.GetAllCustomers(r.Context())
	if err != nil {
		loggers.ErrorLogger.Println("handleGetAllCustomers s.customersSvc.GetAllCustomers error:", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err = jsoner(w, customersArr, statusCode)
	if err != nil {
		loggers.ErrorLogger.Println("handleGetAllCustomers jsoner error:", err)
		return
	}
}
