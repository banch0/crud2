package token

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/banch0/crud2/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidPassword ...
var ErrInvalidPassword = errors.New("invalid password")

// ErrInvalidLogin ...
var ErrInvalidLogin = errors.New("invalid login")

// Service ...
type Service struct {
	pool   *pgxpool.Pool
	secret []byte
}

// RequestDTO ...
type RequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ResponseDTO ...
type ResponseDTO struct {
	Token string `json:"token"`
}

// NewService ...
func NewService(secret jwt.Secret, pool *pgxpool.Pool) *Service {
	return &Service{secret: secret, pool: pool}
}

// Generate ...
func (s *Service) Generate(ctx context.Context, request *RequestDTO) (response ResponseDTO, err error) {
	var pass string
	err = s.pool.QueryRow(ctx, `SELECT password FROM users WHERE login = $1;`, request.Login).Scan(&pass)
	if err != nil {
		log.Println(err)
	}

	header := &jwt.JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(request.Password))
	if err != nil {
		log.Println("Compare hash error: ", err)
		err = ErrInvalidPassword
		return
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = ErrInvalidPassword
		return
	}

	response.Token, err = jwt.Encode(header, jwt.JWTPayload{
		ID:      1,
		Name:    request.Login,
		Expired: time.Now().Add(24 * time.Hour).Unix(),
		Roles:   []string{"ROLE_ADMIN", "ROLE_USER"},
	}, s.secret)

	return
}

// Create ...
func (s *Service) Create(ctx context.Context, request *RequestDTO) (err error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
	}

	_, err = s.pool.Exec(ctx, `INSERT INTO users (login, password, roles) VALUES ($1, $2, $3);`, request.Login, hash, []string{"ROLE_ADMIN", "ROLE_USER"})
	if err != nil {
		log.Println(err)
	}

	return
}
