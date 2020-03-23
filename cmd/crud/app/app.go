package app

import (
	"errors"
	"log"
	"net/http"
	"os"

	agent "github.com/banch0/crud2/cmd/crud/app/api/agents"
	house "github.com/banch0/crud2/cmd/crud/app/api/houses"
	api "github.com/banch0/crud2/cmd/crud/app/api/owners"
	"github.com/banch0/crud2/pkg/crud/services/agents"

	"github.com/banch0/crud2/pkg/crud/services/houses"
	"github.com/banch0/crud2/pkg/crud/services/owners"
	"github.com/banch0/crud2/pkg/jwt"
	"github.com/banch0/crud2/pkg/mux"
	"github.com/banch0/crud2/pkg/token"
	"github.com/banch0/crud2/pkg/user"
	utils "github.com/banch0/crud2/pkg/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ErrorDTO ...
type ErrorDTO struct {
	Errors []string `json:"errors"`
}

// Server ...
type Server struct {
	router    *mux.ExactMux
	pool      *pgxpool.Pool
	secret    jwt.Secret
	tokenSvc  *token.Service
	userSvc   *user.Service
	agentsSvc *agents.Service
	houseSvc  *houses.Service
	ownserSvc *owners.Service
	apiOwner  *api.OwnerServer
	apiAgent  *agent.MainAgentServer
	apiHouse  *house.ServiceHouse
}

// NewServer ...
func NewServer(router *mux.ExactMux,
	pool *pgxpool.Pool,
	secret jwt.Secret,
	tokenSvc *token.Service,
	userSvc *user.Service,
	agentsSvc *agents.Service,
	houseSvc *houses.Service,
	ownserSvc *owners.Service,
	apiOwner *api.OwnerServer,
	apiAgent *agent.MainAgentServer,
	apiHouse *house.ServiceHouse,
) *Server {
	return &Server{
		router:    router,
		pool:      pool,
		secret:    secret,
		tokenSvc:  tokenSvc,
		userSvc:   userSvc,
		agentsSvc: agentsSvc,
		ownserSvc: ownserSvc,
		houseSvc:  houseSvc,
		apiOwner:  apiOwner,
		apiAgent:  apiAgent,
		apiHouse:  apiHouse,
	}
}

// Start ...
func (s *Server) Start() {
	s.InitRoutes()
}

// Stop ...
func (s *Server) Stop() {
	os.Exit(1)
}

// ServerHTTP ...
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleProfile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		response, err := s.userSvc.Profile(request.Context())
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := utils.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_requst"},
			})
			log.Println(err)
			return
		}
		err = utils.WriteJSONBody(writer, &response)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) handleCreateToken() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body token.RequestDTO

		err := utils.ReadJSONBody(request, &body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := utils.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.json_invalid"},
			})
			log.Print(err)
			return
		}

		response, err := s.tokenSvc.Generate(request.Context(), &body)
		if err != nil {
			switch {
			case errors.Is(err, token.ErrInvalidLogin):
				writer.WriteHeader(http.StatusBadRequest)
				err := utils.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.login_mismatch"},
				})
				log.Print(err)
			case errors.Is(err, token.ErrInvalidPassword):
				writer.WriteHeader(http.StatusBadRequest)
				err := utils.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.password_mismatch"},
				})
				log.Print(err)
			default:
				writer.WriteHeader(http.StatusBadRequest)
				err := utils.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.unknown"},
				})
				log.Print(err)
			}
			return
		}

		err = utils.WriteJSONBody(writer, &response)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		user := &token.RequestDTO{}
		err := utils.ReadJSONBody(request, &user)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := utils.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.json_invalid"},
			})
			log.Print(err)
		}

		err = s.tokenSvc.Create(request.Context(), user)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := utils.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.token_create"},
			})
			log.Print(err)
		}

		_, err = writer.Write([]byte(`{"resgister":"sucssesful"}`))
		if err != nil {
			log.Print(err)
		}
	}
}
