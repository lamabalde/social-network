package sqlite

import (
	"errors"
	"os"
	"testing"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

const (
	T_DBPath          = "../testDB.db"
	T_MigratePathTest = "../migrations/sqlite"
	T_MigrateTestDB   = "../migrations/testDB"
)

func TestMigration(t *testing.T) {
	const locTestDB = "test.db"

	err := os.Remove(locTestDB)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	DB, err := OpenDatabase(locTestDB)
	if err != nil {
		t.Fatal("cant open database")
	}
	defer DB.Close()
	dbModel := &queries.DBModel{DB: DB}

	printMigrationVersion(dbModel)

	if err := applyUpMigrations(DB, T_MigratePathTest); err != nil { // this migrates all the way to the latest version
		t.Fatalf("migration failed: %v\n", err)
	}
	printMigrationVersion(dbModel)

	if err := applyUpMigrations(DB, T_MigrateTestDB); err != nil { // this migrates all the way to the latest version
		t.Fatalf("migration failed: %v\n", err)
	}
	printMigrationVersion(dbModel)

	if err := applyDownMigrations(DB, T_MigrateTestDB); err != nil {
		t.Fatalf("migration failed: %v\n", err)
	}
	printMigrationVersion(dbModel)

	if err := applyDownMigrations(DB, T_MigratePathTest); err != nil {
		t.Fatalf("migration failed: %v\n", err)
	}
	printMigrationVersion(dbModel)
}

func TestStepMigration(t *testing.T) {
	const locTestDB = "test.db"

	err := os.Remove(locTestDB)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	DB, err := OpenDatabase(locTestDB)
	if err != nil {
		t.Fatal("cant open database")
	}
	defer DB.Close()
	dbModel := &queries.DBModel{DB: DB}

	if err := applyUpMigrations(DB, T_MigratePathTest); err != nil { // this migrates all the way to the latest version
		t.Fatalf("migration failed: %v\n", err)
	}
	printMigrationVersion(dbModel)

	applySpecificMigrations(DB, T_MigrateTestDB, 5)
	printMigrationVersion(dbModel)

	applySpecificMigrations(DB, T_MigrateTestDB, -6)
	printMigrationVersion(dbModel)
	applySpecificMigrations(DB, T_MigrateTestDB, -8)
	printMigrationVersion(dbModel)

	applySpecificMigrations(DB, T_MigrateTestDB, 2)
	printMigrationVersion(dbModel)
}
