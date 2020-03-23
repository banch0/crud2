package agent

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/banch0/crud2/pkg/crud/models"
	"github.com/banch0/crud2/pkg/crud/services/agents"
	"github.com/banch0/crud2/pkg/mux"
	utils "github.com/banch0/crud2/pkg/utils"
)

// MainAgentServer ...
type MainAgentServer struct {
	agentsSvc *agents.Service
}

// NewMainServer ...
func NewMainServer(agentsSvc *agents.Service) *MainAgentServer {
	return &MainAgentServer{agentsSvc: agentsSvc}
}

// HandleGetAllAgents ...
func (m *MainAgentServer) HandleGetAllAgents(writer http.ResponseWriter, request *http.Request) {
	models, err := m.agentsSvc.GetAllAgents(request.Context())
	if err != nil {
		log.Println("eror", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONBody(writer, &models)
	if err != nil {
		log.Print(err)
	}
}

// HandleAgentByID data validation
func (m *MainAgentServer) HandleAgentByID(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	agent, err := m.agentsSvc.GetAgentByID(request.Context(), id)
	if err != nil {
		log.Println(err)
	}
	err = utils.WriteJSONBody(writer, &agent)
	if err != nil {
		log.Print(err)
	}
}

// HandleAddAgent ...
func (m *MainAgentServer) HandleAddAgent(writer http.ResponseWriter, request *http.Request) {
	agent := &agents.Agents{}
	err := json.NewDecoder(request.Body).Decode(&agent)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.agentsSvc.CreateAgent(request.Context(), *agent)
	if err != nil {
		log.Println(err)
	}

	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success create object"},
	})
}

// HandleUpdateAgent ...
func (m *MainAgentServer) HandleUpdateAgent(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	agent := &agents.Agents{}
	err := json.NewDecoder(request.Body).Decode(&agent)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	agent.ID = id
	err = m.agentsSvc.UpdateAgent(request.Context(), agent)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success updated agent"},
	})
}

// HandleDeleteAgent ...
func (m *MainAgentServer) HandleDeleteAgent(writer http.ResponseWriter, request *http.Request) {
	id, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := m.agentsSvc.DeleteByID(request.Context(), id)
	if err != nil {
		log.Println("Error deleting agent: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSONBody(writer, &models.ConsoleReponse{
		Response: []string{"success removed object"},
	})
}

// AllAgentHouses ...
func (m *MainAgentServer) AllAgentHouses(writer http.ResponseWriter, request *http.Request) {
	models, err := m.agentsSvc.GetAllAgents(request.Context())
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
