package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/banch0/crud2/pkg/crud/models"
	"github.com/banch0/crud2/pkg/crud/services/owners"
	"github.com/banch0/crud2/pkg/mux"
	utils "github.com/banch0/crud2/pkg/utils"
)

// OwnerServer ...
type OwnerServer struct {
	ownerSvc *owners.Service
}

// NewMainServer ...
func NewMainServer(ownerSvc *owners.Service) *OwnerServer {
	return &OwnerServer{ownerSvc: ownerSvc}
}

// HandleGetAllOwners ...
func (m *OwnerServer) HandleGetAllOwners(writer http.ResponseWriter, request *http.Request) {
	models, err := m.ownerSvc.GetOwners(request.Context())
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSONBody(writer, &models)
	if err != nil {
		log.Print(err)
	}
}

// HandleDeleteOwner ...
func (m *OwnerServer) HandleDeleteOwner(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := m.ownerSvc.DeleteOwner(request.Context(), id)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success removed owner"},
	})
}

// HandleGetOwnerByID data validation
func (m *OwnerServer) HandleGetOwnerByID(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	owner, err := m.ownerSvc.GetOwnerByID(request.Context(), id)
	if err != nil {
		log.Println(err)
	}
	err = utils.WriteJSONBody(writer, &owner)
	if err != nil {
		log.Print(err)
	}
}

// HandleAddOwner ...
func (m *OwnerServer) HandleAddOwner(writer http.ResponseWriter, request *http.Request) {
	own := &models.Owners{}
	err := json.NewDecoder(request.Body).Decode(&own)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.ownerSvc.AddNewOwner(request.Context(), *own)
	if err != nil {
		log.Println(err)
	}

	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success create owner"},
	})
}

// HandleUpdateOwner ...
func (m *OwnerServer) HandleUpdateOwner(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	own := &models.Owners{}
	err := json.NewDecoder(request.Body).Decode(&own)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	own.ID = id
	err = m.ownerSvc.UpdateOwner(request.Context(), *own)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success updated owner"},
	})
}
