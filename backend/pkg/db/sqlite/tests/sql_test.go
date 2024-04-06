package tests

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

const (
	T_DBPath          = "../../testDB.db"
	T_MigratePathTest = "../../migrations/sqlite"
	T_MigrateTestDB   = "../../migrations/testDB"
)

func TestCreateDB(t *testing.T) {
	err := os.Remove(T_DBPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	db, err := sqlite.CreateDB(T_DBPath, T_MigratePathTest, T_MigrateTestDB, 1000)
	if err != nil {
		t.Fatal(err)
	}

	defer db.DB.Close()

	fmt.Println("-----GetAllUsers-------")
	uss, err := db.GetAllUsers()
	if err != nil {
		t.Fatal(err)
	}

	for _, us := range uss {
		fmt.Println(us)
	}
	fmt.Println("------------")

	fmt.Println("------AddUser------")

	user := models.User{
		ID:         "4",
		UserName:   "test11",
		Email:      "test1@email",
		DateCreate: time.Date(2023, time.March, 3, 12, 12, 21, 0, time.UTC),
		DateBirth:  time.Date(2002, time.March, 3, 0, 0, 0, 0, time.UTC),
		Gender:     "He",
		FirstName:  "John",
		LastName:   "Test",
	}
	pass, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost) // hash the password
	user.Password = string(pass)

	_, err = db.AddUser(&user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("------------")

	fmt.Println("------GetAllUsers------")
	uss, err = db.GetAllUsers()
	if err != nil {
		t.Fatal(err)
	}

	for _, us := range uss {
		fmt.Println(us)
	}
	fmt.Println("----end-----")
}

func TestAuthenDB(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var sqlconn *sqlite3.SQLiteConn
	err = sqlconn.AuthUserAdd("webuser1", "webuser", false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("----end-----")
}

func getRowNumbersInTable(DB *sql.DB, table, condition string, arguments []any) (int, error) {
	var rowsNum int
	err := DB.QueryRow(`SELECT count(*) FROM `+table+condition, arguments...).Scan(&rowsNum)
	if err != nil {
		return 0, err
	}
	return rowsNum, nil
}
