package db

import (
	"context"
	"errors"
	"homie/models"
	"log"

	"github.com/jackc/pgx/v4"
)

// AddNewCity ...
func AddNewCity(ctx context.Context, cit models.Cities) (err error) {
	sqlQuery := `insert into cities (city) values ($1);`
	_, err = pgPool.Exec(ctx, sqlQuery, &cit.City)
	if err != nil {
		log.Println(err)
	}
	return err
}

// GetCityByID ...
func GetCityByID(ctx context.Context, ID string) (cit models.Cities, err error) {
	sqlQuery := `select id, city from cities where id = $1`
	log.Println(ID)
	err = pgPool.QueryRow(ctx, sqlQuery, ID).Scan(&cit.ID, &cit.City)
	if err != pgx.ErrNoRows {
		err = nil
	}
	return
}

// GetCities ...
func GetCities(ctx context.Context) (cities []*models.Cities, err error) {
	sqlQuery := `select * from cities`
	rows, err := pgPool.Query(ctx, sqlQuery)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		p := new(models.Cities)
		err := rows.Scan(&p.ID, &p.City)
		if err != nil {
			log.Println("Error from db query")
		}
		cities = append(cities, p)
	}

	return cities, err
}

// DeleteCity ...
func DeleteCity(ctx context.Context, ID string) (err error) {
	sqlQuery := `DELETE FROM cities WHERE id = $1`
	res, err := pgPool.Exec(ctx, sqlQuery, ID)
	if err := res.RowsAffected(); err == 0 {
		return errors.New("failed")
	}
	return err
}

// UpdateCity ...
func UpdateCity(ctx context.Context, cit models.Cities) (err error) {
	sqlQuery := `UPDATE cities SET id = $1,
	city = $2 where id = $1`
	result, err := pgPool.Exec(ctx, sqlQuery, cit.ID, cit.City)
	if err != nil {
		log.Println(err)
	}

	if ir := result.RowsAffected(); ir == 0 {
		return errors.New("failed")
	}

	return err
}
