package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

type DBModel struct {
	DB *sql.DB
}

/*
checks if the table exists in DB
*/
func (dbm *DBModel) CheckExistingTable(tableName string) bool {
	row := dbm.DB.QueryRow(`SELECT  name FROM sqlite_master WHERE type='table' AND name = ?`, tableName)

	var tmp string
	err := row.Scan(&tmp)

	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Fatalf("sqlite failed to scan: %v", err)
	}
	return true
}

func (dbm *DBModel) GetMigrationVersion() (models.Migration, error) {
	row := dbm.DB.QueryRow(`SELECT * FROM schema_migrations`)

	var migr models.Migration

	err := row.Scan(&migr.Version, &migr.Dirty)
	if err != nil {
		log.Printf("sqlite failed to scan: %v", err)
		return migr, err
	}

	return migr, nil
}

/*
checks if the value exists in the table's field and returns the number of rows where the value was found
*/
func (dbm *DBModel) checkExisting(table, field, value string) error {
	q := `SELECT ` + field + ` FROM ` + table + ` WHERE ` + field + ` = ?`
	row := dbm.DB.QueryRow(q, value)
	var tmp string
	return row.Scan(&tmp)
}

/*
checks the res and returns error=nil if only 1 row had been affected,
in the other cases returns  ErrNoRecord (for 0 rows), or ErrTooManyRecords (for more than 1)
*/
func (dbm *DBModel) checkUnique(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 1 {
		return nil
	}
	if n == 0 {
		return models.ErrNoRecords
	}
	if n > 1 {
		return models.ErrTooManyRecords
	}
	return errors.New("negative number of rows")
}

func (dbm *DBModel) runTransations(query string, args ...any) error {
	tx, err := dbm.DB.Begin()
	if err != nil {
		return fmt.Errorf("cant start transaction: %w", err)
	}

	// try exec transaction
	_, err = tx.Exec(query, args...)
	if err != nil {
		errRoll := tx.Rollback()
		if errRoll != nil {
			return fmt.Errorf("exec transaction failed: %w, unable to rollback: %w", err, errRoll)
		}
		return fmt.Errorf("exec transaction failed, rolled back: %w", err)
	}

	// if the transaction was a success
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (dbm *DBModel) setFieldStringWhereId(table, field, value, id string) error {
	q := `UPDATE ` + table + ` SET ` + field + `=? WHERE id=?`

	res, err := dbm.DB.Exec(q, value, id)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}