package tests

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wshub"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

func TestGetFollowers(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	tests := []struct {
		userID string
		exp    []*models.UserBase
	}{
		{
			userID: "1",
			exp: []*models.UserBase{
				{ID: "3", UserName: "test2"},
			},
		},
		{
			userID: "2",
			exp: []*models.UserBase{
				{ID: "1", UserName: "no"},
			},
		},
		{
			userID: "3",
			exp: []*models.UserBase{
				{ID: "1", UserName: "no"},
				{ID: "2", UserName: "test1"},
			},
		},
		{
			userID: "4",
			exp:    nil,
		},
	}

	for _, test := range tests {
		res, err := f.GetFollowers(test.userID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(res, test.exp) {
			t.Fatalf("Expected %v, got %v", test.exp, res)
		}
	}
}

func TestGetFollowings(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	tests := []struct {
		userID string
		exp    []models.UserBase
	}{
		{
			userID: "1",
			exp: []models.UserBase{
				{ID: "2", UserName: "test1"},
				{ID: "3", UserName: "test2"},
				//{ID: "4", UserName: "test22"},
			},
		},
		{
			userID: "2",
			exp: []models.UserBase{
				{ID: "3", UserName: "test2"},
			},
		},
		{
			userID: "3",
			exp: []models.UserBase{
				{ID: "1", UserName: "no"},
			},
		},
		{
			userID: "4",
			exp:    nil,
		},
	}

	for _, test := range tests {
		res, err := f.GetFollowing(test.userID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(res, test.exp) {
			t.Fatalf("Expected %v, got %v", test.exp, res)
		}
	}
}

func TestGetFollows(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	tests := []struct {
		userID string
		exp    []models.UserBase
	}{
		{
			userID: "1",
			exp: []models.UserBase{
				{ID: "2", UserName: "test1"},
				{ID: "3", UserName: "test2"},
				//{ID: "4", UserName: "test22"},
			},
		},
		{
			userID: "3",
			exp: []models.UserBase{
				{ID: "1", UserName: "no"},
				{ID: "2", UserName: "test1"},
			},
		},
		{
			userID: "4",
			exp:    nil, // []models.UserBase{
			// 	{ID: "1", UserName: "no"},
			// },

		},
		{
			userID: "5",
			exp:    nil,
		},
	}

	for i, test := range tests {
		res, err := f.GetFollows(test.userID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(res, test.exp) {
			t.Fatalf("Test# %d: Expected %v, got %v", i, test.exp, res)
		}
	}
}

func TestGetFilteredFollowsOrderedByMessagesToGivenUser(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	usersON := wshub.MapID{"1": nil, "2": nil, "3": nil, "4": nil}
	usersONP := wshub.MapID{"2": nil, "3": nil}

	tests := []struct {
		userID string
		online wshub.MapID
		exp    []*models.User
	}{
		{
			userID: "1",
			online: usersON,
			exp: []*models.User{
				{ID: "2", UserName: "test1", LastMessageDate: "2023-11-22 10:59:33.656479916+00:00"},
				{ID: "3", UserName: "test2", LastMessageDate: ""},
				//{ID: "4", UserName: "test22", LastMessageDate: ""},
			},
		},
		{
			userID: "2",
			online: usersON,
			exp: []*models.User{
				{ID: "1", UserName: "no", LastMessageDate: "2023-11-22 10:58:33.656479916+00:00"},
				{ID: "3", UserName: "test2", LastMessageDate: ""},
			},
		},
		{
			userID: "3",
			online: usersON,
			exp: []*models.User{
				{ID: "1", UserName: "no", LastMessageDate: "2023-11-24 10:59:33.656479916+00:00"},
				{ID: "2", UserName: "test1", LastMessageDate: ""},
			},
		},
		{
			userID: "4",
			online: usersON,
			exp:    nil,
		},
		{
			userID: "1",
			online: usersONP,
			exp: []*models.User{
				{ID: "2", UserName: "test1", LastMessageDate: "2023-11-22 10:59:33.656479916+00:00"},
				{ID: "3", UserName: "test2", LastMessageDate: ""},
			},
		},
		{
			userID: "2",
			online: usersONP,
			exp: []*models.User{
				{ID: "3", UserName: "test2", LastMessageDate: ""},
			},
		},
		{
			userID: "3",
			online: usersONP,
			exp: []*models.User{
				{ID: "2", UserName: "test1", LastMessageDate: ""},
			},
		},
		{
			userID: "4",
			online: usersONP,
			exp:    nil,
		},
		{
			userID: "5",
			online: usersON,
			exp:    nil,
		},
	}

	for i, test := range tests {
		res, err := f.GetFilteredFollowsOrderedByMessagesToGivenUser(test.online, test.userID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(test.exp, res) {
			t.Fatalf("Test# %d: Expected %v, got %v", i, test.exp, res)
		}
	}
}

// func TestGetFilteredUsersOrderedByMessageDate(t *testing.T) {
// 	db, err := sqlite.OpenDatabase(DBPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()
// 	f := queries.DBModel{DB: db}

// 	users := chat.MapID{1: nil, 2: nil, 3: nil, 4: nil, 5: nil, 6: nil, 7: nil, 8: nil, 9: nil, 10: nil}
// 	// users :=  make(chat.MapID)

// 	uss, err := f.GetFilteredUsersOrderedByMessagesToGivenUser(users, 4)
// 	if err != nil {
// 		t.Fatalf("error: %v", err)
// 	}
// 	fmt.Println("---------------")
// 	for _, user := range uss {
// 		fmt.Printf("user: = %v, last message: %s\n", user, user.LastMessageDate)
// 	}
// }

func TestAddDeleteFollower(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	numofRows, err := getRowNumbersInTable(f.DB, "followers", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add ----------------------------------------------------------------
	tests := []struct {
		userID        []string
		status        string
		expErrContain string
		exp           int
	}{
		{
			userID:        []string{"2", "3"},
			status:        models.FOLLOW_STATUS_FOLLOWING,
			expErrContain: "constraint",
			exp:           numofRows,
		},

		{
			userID:        []string{"2", "2"},
			status:        models.FOLLOW_STATUS_FOLLOWING,
			expErrContain: "constraint",
			exp:           numofRows,
		},

		{
			userID:        []string{"5", "2"},
			status:        models.FOLLOW_STATUS_FOLLOWING,
			expErrContain: "constraint",
			exp:           numofRows,
		},

		{
			userID:        []string{"3", "2"},
			status:        models.FOLLOW_STATUS_REQUESTED,
			expErrContain: "",
			exp:           numofRows + 1,
		},
	}

	for i, test := range tests {
		err := f.AddFollowing(test.userID[0], test.userID[1], test.status)
		if test.expErrContain == "" && err != nil {
			t.Fatalf("test #%d: error: %v", i, err)
		}
		numofRowsres, err2 := getRowNumbersInTable(f.DB, "followers", "", nil)
		if err2 != nil {
			t.Fatal(err)
		}

		if test.exp != numofRowsres {
			t.Fatalf("test #%d: Expected %v, got %v", i, test.exp, numofRowsres)
		}

		if test.expErrContain != "" {
			fmt.Printf("Error: %v\n", err)
			if err == nil {
				t.Fatalf("test #%d: Expected constraint error, got nil", i)
			}
			if !strings.Contains(err.Error(), test.expErrContain) {
				t.Fatalf("test #%d: Expected constraint error, got %v", i, err.Error())
			}
		}
	}

	// Delete ----------------------------------------------------------------

	tests = []struct {
		userID        []string
		status        string
		expErrContain string
		exp           int
	}{
		{
			userID:        []string{"3", "2"},
			expErrContain: "",
			exp:           numofRows,
		},
	}

	for _, test := range tests {
		err = f.DeleteFollowing(test.userID[0], test.userID[1])
		if err != nil {
			t.Fatal(err)
		}

		numofRowsres, err := getRowNumbersInTable(f.DB, "followers", "", nil)
		if err != nil {
			t.Fatal(err)
		}

		if numofRowsres != test.exp {
			t.Fatalf("Expected %v, got %v", numofRowsres, test.exp)
		}

	}
}
