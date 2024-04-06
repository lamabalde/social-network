package tests

import (
	"fmt"
	"reflect"
	//"forum/controllers/chat"
	"testing"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
	"golang.org/x/crypto/bcrypt"
)

func TestAddUserSession(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("--- add a session to the user 1 ---")
	fmt.Println(f.AddUserSession("1", time.Now().Add(60*time.Second)))
	fmt.Println("--- add a session to the user 10(not existing) ---")
	fmt.Println(f.AddUserSession("10", time.Now().Add(60*time.Second)))
}

func TestDeleteUsersSession(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("--- delete the session ses1 ---")
	fmt.Println("err: ", f.DeleteUsersSession("ses1"))
	fmt.Println("--- delete the session c58f0ec5-8807-490a-98c5-9ad9fdb73f43 ---")
	fmt.Println("err: ", f.DeleteUsersSession("c58f0ec5-8807-490a-98c5-9ad9fdb73f43"))
}

func TestInsertUser(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("--- insert a user -ussertest1-  ---")
	user := models.User{
		UserName:   "ussertest1",
		Email:      "emailt1@email",
		DateCreate: time.Now(),
		DateBirth:  time.Date(2002, time.March, 12, 0, 0, 0, 0, time.UTC),
		Gender:     "He",
		FirstName:  "John",
		LastName:   "First",
	}
	pass, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost) // hash the password
	user.Password = string(pass)
	err = f.InsertUser(&user)
	fmt.Printf("err= %v\n", err)

	fmt.Println("--- add a user -usserAdd-  ---")
	user = models.User{
		UserName:   "usserAdd",
		Email:      "emailtAdd@email",
		DateCreate: time.Now(),
		DateBirth:  time.Date(2002, time.March, 3, 0, 0, 0, 0, time.UTC),
		Gender:     "He",
		FirstName:  "John",
		LastName:   "second",
	}
	pass, _ = bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost) // hash the password
	user.Password = string(pass)

	_, err = f.AddUser(&user)
	fmt.Printf(" err= %v\n", err)
}

func TestGetUserByUUID(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("---GetUserBySession(45f20c10-7d59-49cf-b512-fd27704f8436)---")
	u, uuid, exp, err := f.GetUserBySession("45f20c10-7d59-49cf-b512-fd27704f8436")
	fmt.Printf("user= %v, uuid=%v, exp: %v\nerr= %v\n", u, uuid, exp, err)

	fmt.Println("---no that session- d8ce41bc-a504-4c4d-9285-c560a4b--")
	u, uuid, exp, err = f.GetUserBySession("d8ce41bc-a504-4c4d-9285-c560a4b")
	fmt.Printf("user= %v, uuid=%v, exp: %v\nerr= %v\n", u, uuid, exp, err)
}

func TestGetUserByID(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("---id 2---")
	u, err := f.GetUserByID("2")
	fmt.Printf("user= %s, \nerr= %v\n", u.StringFull(), err)

	fmt.Println("---id 4---")
	u, err = f.GetUserByID("4")
	fmt.Printf("user= %s, \nerr= %v\n", u.StringFull(), err)
}

func TestGetUserByName(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("---name test1---")
	u, err := f.GetUserByName("test1")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)

	fmt.Println("---name test4---")
	u, err = f.GetUserByName("test4")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)
}

func TestGetUserByEmail(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("---Email test1---")
	u, err := f.GetUserByEmail("test1@forum")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)

	fmt.Println("---Email test4---")
	u, err = f.GetUserByEmail("test@ff.f4")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)
}

func TestGetUserByNameOrEmail(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("---Email test1---")
	u, err := f.GetUserByNameOrEmail("test1@forum")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)

	fmt.Println("---name test1---")
	u, err = f.GetUserByNameOrEmail("test1")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)

	fmt.Println("---name noname--")
	u, err = f.GetUserByNameOrEmail("noname")
	fmt.Printf("user= %v, \nerr= %v\n", u, err)
}

func TestChangeUsersProfileType(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	tests := []struct {
		userID string
		prtype int
	}{
		{
			userID: "1",
			prtype: 1,
		},
		{
			userID: "2",
			prtype: 1,
		},
		{
			userID: "2",
			prtype: 1,
		},
		{
			userID: "1",
			prtype: 0,
		},
		{
			userID: "2",
			prtype: 0,
		},
	}

	for _, test := range tests {
		err := f.SetUsersProfileType(test.userID, test.prtype)
		if err != nil {
			t.Fatal(err)
		}

		res := -1
		err = f.DB.QueryRow("SELECT profiletype FROM users WHERE id=?", test.userID).Scan(&res)
		if err != nil {
			t.Fatal(err)
		}

		if test.prtype != res {
			t.Fatalf("Expected %v, got %v", test.prtype, res)
		}
	}
}

func TestGetUsers(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	tests := []struct {
		groupID     string
		searchQuery string
		exp         []*models.UserBase
	}{
		{
			groupID: "",
			exp: []*models.UserBase{
				{ID: "1", UserName: "no"},
				{ID: "2", UserName: "test1"},
				{ID: "3", UserName: "test2"},
				{ID: "4", UserName: "test22"},
			},
		},
		{
			groupID:     "3",
			searchQuery: "te",
			exp: []*models.UserBase{
				{ID: "2", UserName: "test1"},
				{ID: "3", UserName: "test2"},
				{ID: "4", UserName: "test22"},
			},
		},
		{
			groupID:     "3",
			searchQuery: "2",
			exp: []*models.UserBase{
				{ID: "3", UserName: "test2"},
				{ID: "4", UserName: "test22"},
			},
		},
		{
			groupID:     "2",
			searchQuery: "te",
			exp: []*models.UserBase{
				{ID: "4", UserName: "test22"},
			},
		},
		{
			groupID:     "2",
			searchQuery: "1",
			exp:         nil,
		},
		{
			groupID:     "3",
			searchQuery: "1",
			exp: []*models.UserBase{
				{ID: "2", UserName: "test1"},
			},
		},
		{
			groupID:     "3",
			searchQuery: "f",
			exp:         nil,
		},
	}

	for i, test := range tests {
		var err error
		var res []*models.UserBase
		if test.groupID == ""{
			res, err = f.GetAllUsers()	
		}else {
			res, err = f.GetUsersByPartialNameNotInGroup(test.searchQuery,test.groupID)
		}
		if err != nil {
			t.Fatalf("test #%d,Error: %#v",i,err)
		}

		if !reflect.DeepEqual(test.exp,res) {
			t.Fatalf("Test#%d. Expected %v, got %v",i, test.exp, res)
		}
	}
}
