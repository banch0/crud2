package agents

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Agents ...
type Agents struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Login      string   `json:"login"`
	Password   string   `json:"password"`
	Surname    string   `json:"lastname"`
	Email      string   `json:"email"`
	Individual string   `json:"individual"`
	Logo       string   `json:"logo"`
	Phone      string   `json:"phone"`
	Roles      []string `json:"roles"`
}

// Service ...
type Service struct {
	pool *pgxpool.Pool
}

// NewService ...
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

var createDB = `
CREATE TABLE IF NOT EXISTS agents (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	login TEXT UNIQUE,
	password TEXT,
	roles TEXT[],
	removed BOOLEAN DEFAULT FALSE,
	name VARCHAR(50) NOT NULL,
	lastname VARCHAR(50) NOT NULL,
	phone VARCHAR(20),
	individual BOOLEAN DEFAULT FALSE,
	craetedAt TIMESTAMP DEFAULT now(),
	logo TEXT DEFAULT '',
	others jsonb DEFAULT '{}',
	PRIMARY KEY(id)
);
`

// GetAllAgents ...
func (service *Service) GetAllAgents(ctx context.Context) (models []Agents, err error) {
	rows, err := service.pool.Query(context.Background(), `
	SELECT id, name, lastname, phone FROM agents WHERE removed = FALSE;
	`)
	if err != nil {
		return nil, fmt.Errorf("can't get agents from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		agent := Agents{}
		err = rows.Scan(&agent.ID, &agent.Name, &agent.Surname, &agent.Phone)
		if err != nil {
			return nil, fmt.Errorf("can't get agents from db: %w", err)
		}
		models = append(models, agent)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get agents from db: %w", err)
	}
	return models, nil
}

// GetAgentByID ...
func (service *Service) GetAgentByID(ctx context.Context, id string) (agent Agents, err error) {
	sqlQuery := `SELECT id, name, lastname, login, email, phone FROM agents WHERE id = $1`
	err = service.pool.QueryRow(ctx, sqlQuery, id).Scan(&agent.ID, &agent.Name, &agent.Surname, &agent.Login, &agent.Email, &agent.Phone)
	if err != pgx.ErrNoRows {
		err = nil
	}
	return
}

// CreateAgent ...
func (service *Service) CreateAgent(ctx context.Context, agent Agents) (err error) {
	sqlQuery := `INSERT INTO agents (name, lastname, phone, password) VALUES ($1, $2, $3, $4);`
	_, err = service.pool.Exec(ctx, sqlQuery, &agent.Name, &agent.Surname, &agent.Phone, &agent.Password)
	if err != nil {
		log.Println(err)
	}
	return err
}

// DeleteByID ...
func (service *Service) DeleteByID(ctx context.Context, id string) (err error) {
	sqlQuery := `UPDATE agents SET removed = TRUE WHERE id = $1`
	res, err := service.pool.Exec(ctx, sqlQuery, id)
	if err := res.RowsAffected(); err == 0 {
		return errors.New("failed")
	}
	return err
}

// UpdateAgent ...
func (service *Service) UpdateAgent(ctx context.Context, agent *Agents) (err error) {
	sqlQuery := `UPDATE agents SET name = $2, login = $5, email = $6, lastname = $3, phone = $4 WHERE id = $1`
	result, err := service.pool.Exec(ctx, sqlQuery, &agent.ID, &agent.Name, &agent.Surname, &agent.Phone, &agent.Login, &agent.Email)
	if err != nil {
		log.Println("UpdateOwner error:", err)
	}

	if ir := result.RowsAffected(); ir == 0 {
		return errors.New("failed")
	}

	return err
}
