package app

import (
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	err = r.ParseMultipartForm(64 << 40)
	if err != nil {
		log.Println("handleSaveMedicine r.ParseMultipartForm error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item.ID, err = strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		log.Println("handleSaveMedicine r.FormValue(id) parsing error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item.Name = r.FormValue("name")
	item.Manafacturer = r.FormValue("manafacturer")
	item.Description = r.FormValue("description")
	item.Components = strings.Split(r.FormValue("components"), ", ")
	item.Recipe_needed, err = strconv.ParseBool(r.FormValue("recipe_needed"))
	if err != nil {
		log.Println("handleSaveMedicine r.FormValue(recipe_needed) parsing error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item.Price, err = strconv.Atoi(r.FormValue("price"))
	if err != nil {
		log.Println("handleSaveMedicine r.FormValue(price) parsing error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item.Qty, err = strconv.Atoi(r.FormValue("qty"))
	if err != nil {
		log.Println("handleSaveMedicine r.FormValue(qty) parsing error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item.PharmacyName = r.FormValue("pharmacy_name")
	item.Active, err = strconv.ParseBool(r.FormValue("active"))
	if err != nil {
		log.Println("handleSaveMedicine r.FormValue(active) parsing error:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	imageExt := filepath.Ext(r.FormValue("image"))
	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		var name = strings.Split(fileHeader.Filename, ".")
		imageExt = name[len(name)-1]
	}
	item.Image = imageExt
	item.File = item.Name + "." + imageExt
	dir := strconv.FormatInt(time.Now().Unix(), 64)
	loadFile(file, dir, "../images/", item.File)

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

func loadFile(file multipart.File, dir string, path string, namefile string) error {
	err := os.MkdirAll(path+dir, 0777)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("not readble data")
	}

	err = ioutil.WriteFile(path+dir+"/"+namefile, data, 0666)

	if err != nil {
		return errors.New("not saved from folder ")
	}
	return nil
}
