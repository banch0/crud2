package house

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/banch0/crud2/pkg/crud/models"
	"github.com/banch0/crud2/pkg/crud/services/houses"
	"github.com/banch0/crud2/pkg/mux"
	utils "github.com/banch0/crud2/pkg/utils"
)

// ServiceHouse ...
type ServiceHouse struct {
	houseSvc *houses.Service
}

// NewHouseServer ...
func NewHouseServer(houseSvc *houses.Service) *ServiceHouse {
	return &ServiceHouse{houseSvc: houseSvc}
}

// GetAllHouses ...
func (m *ServiceHouse) GetAllHouses(writer http.ResponseWriter, request *http.Request) {
	houses, err := m.houseSvc.GetAllHouses(request.Context())
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSONBody(writer, &houses)
	if err != nil {
		log.Print(err)
	}
}

// GetHouseByID ...
func (m *ServiceHouse) GetHouseByID(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	house, err := m.houseSvc.GetHouseByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONBody(writer, &house)
	if err != nil {
		log.Print(err)
	}
}

// HousesByCategories ...
func (m *ServiceHouse) HousesByCategories(writer http.ResponseWriter, request *http.Request) {
	category, ok := mux.FromContext(request.Context(), "category")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	houses, err := m.houseSvc.GetHousesByCategory(request.Context(), category)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSONBody(writer, &houses)
	if err != nil {
		log.Print(err)
	}
}

// AddNewHouse ...
func (m *ServiceHouse) AddNewHouse(writer http.ResponseWriter, request *http.Request) {
	house := &models.Houses{}
	err := json.NewDecoder(request.Body).Decode(&house)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = m.houseSvc.CreateNewHouseDB(request.Context(), house)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success adding object"},
	})
}

// UpdateHouse ...
func (m *ServiceHouse) UpdateHouse(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	house := &models.Houses{}
	err := json.NewDecoder(request.Body).Decode(&house)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = m.houseSvc.UpdateHouse(request.Context(), id, *house)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success updated object"},
	})
}

// RemoveHouse ...
func (m *ServiceHouse) RemoveHouse(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := m.houseSvc.DeleteHouse(request.Context(), id)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success removed object"},
	})
}
