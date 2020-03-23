package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/banch0/crud2/pkg/mux/middleware/jwt"
)

// ErrNotToken ...
var ErrNotToken = errors.New("token not found")

// Service ...
type Service struct{}

// NewService ...
func NewService() *Service {
	return &Service{}
}

// ResponseDTO ...
type ResponseDTO struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// Profile ...
func (s *Service) Profile(ctx context.Context) (response ResponseDTO, err error) {
	auth := jwt.FromContext(ctx)
	body, err := json.Marshal(auth)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
	}

	return ResponseDTO{
		ID:     response.ID,
		Name:   response.Name,
		Avatar: "https://i.pravatar.cc/50",
	}, nil
}
