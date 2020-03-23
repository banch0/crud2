package app

import (
	"reflect"

	"github.com/banch0/crud2/pkg/auth"
	jwtcore "github.com/banch0/crud2/pkg/jwt"
	"github.com/banch0/crud2/pkg/mux/middleware/authorized"
	"github.com/banch0/crud2/pkg/mux/middleware/jwt"
	"github.com/banch0/crud2/pkg/mux/middleware/logger"
)

// InitRoutes ...
func (s *Server) InitRoutes() {
	s.router.POST("/api/register", s.handleRegister(), logger.Logger("REGISTER"))
	s.router.POST("/api/tokens", s.handleCreateToken(), logger.Logger("LOGIN"))
	s.router.GET("/api/users/me", s.handleProfile(),
		auth.Auth(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*jwtcore.JWTPayload)(nil)).Elem(), s.secret),
		logger.Logger("AUTHORIZED"),
	)

	s.router.GET(
		"/api/admin",
		s.handleProfile(),
		authorized.Authorized([]string{"ROLE_ADMIN"}, jwt.FromContext),
		auth.Auth(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*jwtcore.JWTPayload)(nil)).Elem(), s.secret),
		logger.Logger("ADMIN"),
	)

	// House
	s.router.POST("/api/houses", s.apiHouse.GetAllHouses)
	s.router.GET("/api/house/{id}", s.apiHouse.GetHouseByID)
	s.router.GET("/api/house/category/{category}", s.apiHouse.HousesByCategories)
	s.router.POST("/api/house", s.apiHouse.AddNewHouse)
	s.router.PUT("/api/house/{id}", s.apiHouse.UpdateHouse)
	s.router.DELETE("/api/house/{id}", s.apiHouse.RemoveHouse)

	// Agents
	s.router.POST("/api/agents", s.apiAgent.HandleGetAllAgents)
	s.router.GET("/api/agent/{id}", s.apiAgent.HandleAgentByID)
	s.router.POST("/api/agent/{id}", s.apiAgent.HandleAddAgent)
	s.router.PUT("/api/agent/{id}", s.apiAgent.HandleUpdateAgent)
	s.router.DELETE("/api/agent/{id}", s.apiAgent.HandleDeleteAgent)
	s.router.GET("/api/agent/{id}", s.apiAgent.AllAgentHouses)

	// Owner
	s.router.GET("/api/owners", s.apiOwner.HandleGetAllOwners)
	s.router.GET("/api/owner/{id}", s.apiOwner.HandleGetOwnerByID)
	s.router.DELETE("/api/owner/{id}", s.apiOwner.HandleDeleteOwner)
	s.router.PUT("/api/owner/{id}", s.apiOwner.HandleUpdateOwner)
	s.router.POST("/api/owner/{id}", s.apiOwner.HandleAddOwner)

}
