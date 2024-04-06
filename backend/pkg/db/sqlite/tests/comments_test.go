package tests

import (
	"fmt"
	"testing"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

func TestInsertComment(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	f := queries.DBModel{DB: db}

	comment, err := f.AddComment("2", "comment 2 to post 2", []string{}, "1")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("--comment--: %v", comment)
}
