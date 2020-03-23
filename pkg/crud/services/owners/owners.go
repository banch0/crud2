package owners

import (
	"context"
	"errors"
	"log"

	"github.com/banch0/crud2/pkg/crud/models"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Service ...
type Service struct {
	pool *pgxpool.Pool
}

// NewService ...
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// AddNewOwner ...
func (service *Service) AddNewOwner(ctx context.Context, own models.Owners) (err error) {
	sqlQuery := `INSERT INTO owners (name, lastname, phone, email, password) VALUES ($1, $2, $3, $4, $5);`
	_, err = service.pool.Exec(ctx, sqlQuery, &own.Name, &own.Lastname, &own.Phone, &own.Email, &own.Password)
	if err != nil {
		log.Println(err)
	}
	return err
}

// GetOwnerByID ...
func (service *Service) GetOwnerByID(ctx context.Context, ID string) (own models.Owners, err error) {
	sqlQuery := `SELECT id, name, lastname, phone FROM owners WHERE id = $1`
	err = service.pool.QueryRow(ctx, sqlQuery, ID).Scan(&own.ID, &own.Name, &own.Lastname, &own.Phone)
	if err != pgx.ErrNoRows {
		err = nil
	}
	return
}

// GetOwners ...
func (service *Service) GetOwners(ctx context.Context) (own []*models.Owners, err error) {
	sqlQuery := `SELECT id, name, lastname, phone FROM owners;`
	rows, err := service.pool.Query(ctx, sqlQuery)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		p := new(models.Owners)
		err := rows.Scan(&p.ID, &p.Name, &p.Lastname, &p.Phone)
		if err != nil {
			log.Println("Scan error: ", err)
		}
		own = append(own, p)
	}

	return own, err
}

// DeleteOwner ...
func (service *Service) DeleteOwner(ctx context.Context, ID string) (err error) {
	sqlQuery := `UPDATE owners SET removed = TRUE WHERE id = $1`
	res, err := service.pool.Exec(ctx, sqlQuery, ID)
	if err := res.RowsAffected(); err == 0 {
		return errors.New("failed delete owner")
	}
	return err
}

// UpdateOwner ...
func (service *Service) UpdateOwner(ctx context.Context, own models.Owners) (err error) {
	sqlQuery := `UPDATE owners SET name = $2, lastname = $3, phone = $4 WHERE id = $1;`
	result, err := service.pool.Exec(ctx, sqlQuery, own.ID, own.Name, own.Lastname, own.Phone)
	if err != nil {
		log.Println("UpdateOwner error:", err)
	}

	if ir := result.RowsAffected(); ir == 0 {
		return errors.New("failed update owner")
	}

	return err
}
