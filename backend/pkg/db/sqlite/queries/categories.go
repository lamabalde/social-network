package queries

import (
	"database/sql"
	"errors"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

/*
inserts a new comment into DB, returns an ID for the comment
*/
func (dbm *DBModel) GetCategories() ([]*models.Category, error) {
	q := `SELECT id, name FROM categories ORDER BY name`
	rows, err := dbm.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parsing the query's result
	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (dbm *DBModel) GetCategoryByID(id string) (*models.Category, error) {
	q := `SELECT id, name FROM categories WHERE id=?`
	category := models.Category{}
	row := dbm.DB.QueryRow(q, id)
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}
	return &category, nil
}
