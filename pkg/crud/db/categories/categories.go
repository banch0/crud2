package db

import (
	"context"
	"errors"
	"homie/models"
	"log"

	"github.com/jackc/pgx"
)

// CreateNewCategory ...
func CreateNewCategory(ctx context.Context, cat models.Category) (err error) {
	sqlQuery := `insert into category (id, title) values ($1, $2);`
	_, err = pool.Exec(ctx, sqlQuery, &cat.ID, &cat.Title)
	if err != nil {
		log.Println(err)
	}
	return err
}

// GetCategoryByID ...
func GetCategoryByID(ctx context.Context, ID string) (categ models.Category, err error) {
	sqlQuery := `select id, title from category where id = $1`
	log.Println(ID)
	err = pool.QueryRow(ctx, sqlQuery, ID).Scan(&categ.ID, &categ.Title)
	if err != pgx.ErrNoRows {
		err = nil
	}
	return
}

// GetCategory ...
func GetCategory(ctx context.Context) (categories []*models.Category, err error) {
	sqlQuery := `select * from category`
	rows, err := pool.Query(ctx, sqlQuery)
	if err != pgx.ErrNoRows {
		err = nil
	}

	for rows.Next() {
		p := new(models.Category)
		err := rows.Scan(&p.ID, &p.Title)
		if err != nil {
			log.Println("Error from db query")
		}
		categories = append(categories, p)
	}

	return categories, err
}

// DeleteCategory ...
func DeleteCategory(ctx context.Context, ID string) (err error) {
	sqlQuery := `DELETE FROM category WHERE id = $1`
	res, err := pool.Exec(ctx, sqlQuery, ID)
	if err := res.RowsAffected(); err == 0 {
		return errors.New("failed")
	}
	return err
}

// UpdateCategory ...
func UpdateCategory(ctx context.Context, categ models.Category) (err error) {
	sqlQuery := `UPDATE category SET id = $1,
	title = $2 where id = $1`
	result, err := pool.Exec(ctx, sqlQuery, categ.ID, categ.Title)
	if err != nil {
		log.Println(err)
	}

	if ir := result.RowsAffected(); ir == 0 {
		return errors.New("failed")
	}

	return err
}
