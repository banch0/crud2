package houses

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

// CreateNewHouseDB ...
func (s *Service) CreateNewHouseDB(ctx context.Context, h *models.Houses) (err error) {
	sqlQuery := `insert into houses (title, price, description) 
	values ($1, $2, $3);`
	_, err = s.pool.Exec(ctx, sqlQuery, h.Title, h.Price, h.Description)
	if err != nil {
		log.Println(err)
	}
	return
}

// GetHouseByID ...
func (s *Service) GetHouseByID(ctx context.Context, ID string) (hs models.Houses, err error) {
	sqlQuery := `SELECT id, title, price, description FROM houses WHERE id = $1;`
	err = s.pool.QueryRow(ctx, sqlQuery, ID).Scan(&hs.ID, &hs.Title, &hs.Price, &hs.Description)
	if err != nil {
		log.Println(err)
		err = nil
	}
	return
}

// GetAllHouses ...
func (s *Service) GetAllHouses(ctx context.Context) (houses []*models.Houses, err error) {
	sqlQuery := `select id, title, price, description from houses`
	rows, err := s.pool.Query(ctx, sqlQuery)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		house := new(models.Houses)
		err := rows.Scan(&house.ID, &house.Title, &house.Price, &house.Description)
		if err != nil {
			log.Println("Error from db query1")
		}
		houses = append(houses, house)
	}

	return houses, err
}

// DeleteHouse ...
func (s *Service) DeleteHouse(ctx context.Context, ID string) (err error) {
	sqlQuery := `UPDATE houses SET remove = TRUE WHERE id = $1`
	res, err := s.pool.Exec(ctx, sqlQuery, ID)
	if err := res.RowsAffected(); err == 0 {
		return errors.New("failed update house")
	}
	return err
}

// UpdateHouse ...
func (s *Service) UpdateHouse(ctx context.Context, id string, hs models.Houses) (err error) {
	sqlQuery := `UPDATE houses SET title = $2, price = $3, description = $4 where id = $1`
	result, err := s.pool.Exec(ctx, sqlQuery, id, hs.Title, hs.Price, hs.Description)
	if err != nil {
		log.Println(err)
	}

	if ir := result.RowsAffected(); ir == 0 {
		return errors.New("failed")
	}

	return err
}

// BelongToAgent ...
func (s *Service) BelongToAgent(ctx context.Context, agent string) (hs []*models.Houses, err error) {
	sqlQuery := `SELECT id, title, price, description, created_at FROM houses WHERE owner_id = $1`

	rows, err := s.pool.Query(ctx, sqlQuery, agent)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		q := new(models.Houses)
		err := rows.Scan(&q.ID, &q.Title, &q.Price, &q.Description, &q.CreatedAt)
		if err != nil {
			log.Fatal(err)
			return hs, err
		}
		log.Println(q)
		hs = append(hs, q)
	}
	return hs, err
}

// GetHousesByCity ...
func (s *Service) GetHousesByCity(ctx context.Context, city string) (hs []*models.Houses, err error) {
	sqlQuery := `SELECT id, title, price, description, created_at FROM houses WHERE city = $1`

	rows, err := s.pool.Query(ctx, sqlQuery, city)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		q := new(models.Houses)
		err := rows.Scan(&q.ID, &q.Title, &q.Price, &q.Description, &q.CreatedAt)
		if err != nil {
			log.Fatal(err)
			return hs, err
		}
		log.Println(q)
		hs = append(hs, q)
	}
	return hs, err
}

// GetHousesByRayons ...
func (s *Service) GetHousesByRayons(ctx context.Context, rayon string) (hs []*models.Houses, err error) {
	sqlQuery := `SELECT id, title, price, description, created_at FROM houses WHERE rayon = $1`

	rows, err := s.pool.Query(ctx, sqlQuery, rayon)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		q := new(models.Houses)
		err := rows.Scan(&q.ID, &q.Title, &q.Price, &q.Description, &q.CreatedAt)
		if err != nil {
			log.Fatal(err)
			return hs, err
		}
		log.Println(q)
		hs = append(hs, q)
	}

	return hs, err
}

// GetHousesByCategory ...
func (s *Service) GetHousesByCategory(ctx context.Context, category string) (houses []models.Houses, err error) {
	sqlQuery := `SELECT id, title, price, description, created_at FROM houses WHERE category = $1`

	rows, err := s.pool.Query(ctx, sqlQuery, category)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		house := new(models.Houses)
		err := rows.Scan(&house.ID, &house.Title, &house.Price, &house.Description, &house.CreatedAt)
		if err != nil {
			log.Fatal(err)
			return houses, err
		}
		houses = append(houses, *house)
	}

	return houses, err
}
