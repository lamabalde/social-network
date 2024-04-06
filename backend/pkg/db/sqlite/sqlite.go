package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFileName    = "pkg/db/social-network.db"
	TestDB        = "pkg/db/testDB.db"
	MigratePath   = "pkg/db/migrations/sqlite"
	MigrateTestDB = "pkg/db/migrations/testDB"
)

func InitDB(test bool, versionDB int) (*queries.DBModel, error) { // have go-migrate run all the up migrations
	var err error
	if test {
		return CreateTestDB()
	}

	DB, err := OpenDatabase(DBFileName)
	if err != nil {
		return nil, err
	}

	dbModel := &queries.DBModel{DB: DB}
	if versionDB != 0 {
		return migrateToVersion(dbModel, versionDB)
	}

	applyUpMigrations(dbModel.DB, MigratePath)

	return dbModel, nil
}

func OpenDatabase(fileName string) (*sql.DB, error) {
	// init pull (not connection)
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", fileName))
	if err != nil {
		return nil, err
	}

	// check connection (create and check)
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func deleteTestDataFromDB(db *sql.DB) error {
	return applyDownMigrations(db, MigrateTestDB)
}

func dropDB(db *sql.DB) error {
	return applyDownMigrations(db, MigratePath)
}

func applyUpMigrations(db *sql.DB, path string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("error opening driver: %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("error creating instance: %v", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrating database: %v", err)
	}
	return nil
}

func applyDownMigrations(db *sql.DB, path string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("error opening driver: %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("error creating instance: %v", err)
	}

	if err = migration.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrating database: %v", err)
	}
	return nil
}

func applySpecificMigrations(db *sql.DB, path string, amount int) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("error opening driver: %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("error creating instance: %v", err)
	}

	if err = migration.Steps(amount); err != nil && err != migrate.ErrNoChange { // if there is no change, then it continues starting the DB
		return fmt.Errorf("error migrating database: %v", err)
	}

	return nil
}

func handleErrAndCloseDB(db *sql.DB, operation string, err error) error {
	errClose := db.Close()
	if errClose != nil {
		return fmt.Errorf("'%s' failed: %w, unable to close DB: %w", operation, err, errClose)
	}
	return fmt.Errorf("DB was closed cause the '%s' failed: %w", operation, err)
}

/*
fills in the DB with data from the given file
*/
func FillInTestDB(db *queries.DBModel, path string) error {
	ForceMigrate(db, 1000)
	return applyUpMigrations(db.DB, path)
}

func ForceMigrate(db *queries.DBModel, version int) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("filling DB: error opening driver: %v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+MigratePath,
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("filling DB: error creating instance: %v", err)
	}

	err = migration.Force(version)
	if err != nil {
		return fmt.Errorf("filling DB: error forcing the version of database: %v", err)
	}
	return nil
}

func CreateDB(fileDB, createMigrationPath, dataMigrationPath string, dataVersion int) (*queries.DBModel, error) {
	var err error
	DB, err := OpenDatabase(fileDB)
	if err != nil {
		return nil, fmt.Errorf("open database %s failed: %v", fileDB, err)
	}

	dbModel := &queries.DBModel{DB: DB}
	if err := applyUpMigrations(dbModel.DB, createMigrationPath); err != nil { // this migrates all the way to the latest version
		return nil, fmt.Errorf("apply up migration from %v failed: %v", createMigrationPath, err)
	}

	err = FillInTestDB(dbModel, dataMigrationPath)
	return dbModel, err
}

func CreateTestDB() (*queries.DBModel, error) {
	DB, err := OpenDatabase(TestDB)
	if err != nil {
		return nil, err
	}

	dbModel := &queries.DBModel{DB: DB}

	ok := dbModel.CheckExistingTable("schema_migrations")
	if ok {
		migr, err := dbModel.GetMigrationVersion()
		if err == nil && migr.Version > 1000 {
			if err := applyUpMigrations(dbModel.DB, MigrateTestDB); err != nil { // this migrates all the way to the latest version
				return nil, fmt.Errorf("apply up migration from %v failed: %v", MigratePath, err)
			}
			return dbModel, nil
		}
	}

	if err := applyUpMigrations(dbModel.DB, MigratePath); err != nil { // this migrates all the way to the latest version
		return nil, fmt.Errorf("apply up migration from %v failed: %v", MigratePath, err)
	}

	err = FillInTestDB(dbModel, MigrateTestDB)
	return dbModel, err
}

func printMigrationVersion(dbModel *queries.DBModel) (models.Migration, error) {
	migr, err := dbModel.GetMigrationVersion()
	if err != nil {
		fmt.Printf("cant get migration version: %v\n", err)
		return migr, err
	}
	fmt.Printf("migration version: %v\n", migr)
	return migr, err
}

func migrateToVersion(dbModel *queries.DBModel, versionDB int) (*queries.DBModel, error) {
	currentVersion, err := printMigrationVersion(dbModel)
	if err != nil {
		return nil, err
	}

	if currentVersion.Dirty == 1 {
		return nil, errors.New("the migration version is dirty")
	}

	if versionDB >= 1000 && currentVersion.Version < 1000 {
		applyUpMigrations(dbModel.DB, MigratePath)
		printMigrationVersion(dbModel)
		ForceMigrate(dbModel, 1000)
		currentVersion, _ = printMigrationVersion(dbModel)
	}
	if versionDB < 1000 && currentVersion.Version >= 1000 {
		applyDownMigrations(dbModel.DB, MigrateTestDB)
		printMigrationVersion(dbModel)
		vers, err := getLastMigrationVersionInDir(MigratePath)
		if err != nil {
			return nil, err
		}
		ForceMigrate(dbModel, vers)
		currentVersion, _ = printMigrationVersion(dbModel)
	}

	migrPath := MigratePath
	if versionDB >= 1000 {
		migrPath = MigrateTestDB
	}
	applySpecificMigrations(dbModel.DB, migrPath, versionDB-currentVersion.Version)
	printMigrationVersion(dbModel)

	return dbModel, nil
}

func getLastMigrationVersionInDir(dir string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	numbers := make([]int, len(files))

	for i, file := range files {

		numbers[i], err = strconv.Atoi(strings.Split(file.Name(), "_")[0])
		if err != nil {
			return 0, err
		}
	}

	sort.Ints(numbers)
	return numbers[len(numbers)-1], nil
}
