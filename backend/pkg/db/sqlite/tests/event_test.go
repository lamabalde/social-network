package tests

import (
	"fmt"
	"strings"
	"testing"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

func TestAddEventMembers(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	numofRows, err := getRowNumbersInTable(f.DB, "group_event_members", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		eventID, userID string
		mark            int
		expErrContain   string
		exp             int
	}{
		{
			eventID:       "2",
			userID:        "2",
			mark:          1,
			expErrContain: "",
			exp:           numofRows + 1,
		},
		{
			eventID:       "4",
			userID:        "3",
			mark:          1,
			expErrContain: models.ErrNoRecords.Error(),
			exp:           numofRows,
		},
		{
			eventID:       "3",
			userID:        "1",
			mark:          1,
			expErrContain: models.ErrNoRecords.Error(),
			exp:           numofRows,
		},
	}
	for i, test := range tests {
		id, err := f.AddEventMember(test.eventID, test.userID, test.mark)
		if test.expErrContain == "" && err != nil {
			t.Fatalf("test #%d: error: %v", i, err)
		}
		numofRowsres, err2 := getRowNumbersInTable(f.DB, "group_event_members", "", nil)
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

		if err == nil {
			_, err2 := f.DB.Exec("DELETE FROM group_event_members WHERE id=?", id)
			if err2 != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestChangeEventMemberOption(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	tests := []struct {
		eventID, userID     string
		mark                int
		expErrContain       string
		exp                 int
		sel_group_membersID int
	}{
		{
			eventID:             "3",
			userID:              "3",
			mark:                1,
			expErrContain:       "",
			sel_group_membersID:  2,
			exp:                 1,
		},
		{
			eventID:             "3",
			userID:              "3",
			mark:                0,
			expErrContain:       "",
			sel_group_membersID: 2,
			exp:                 0,
		},
		{
			eventID:       "2",
			userID:        "2",
			mark:          1,
			expErrContain: models.ErrNoRecords.Error(),
		},
		{
			eventID:       "4",
			userID:        "3",
			mark:          1,
			expErrContain: models.ErrNoRecords.Error(),
		},
		{
			eventID:       "3",
			userID:        "1",
			mark:          1,
			expErrContain: models.ErrNoRecords.Error(),
		},
	}
	for i, test := range tests {
		err := f.ChangeEventMemberOption(test.eventID, test.userID, test.mark)
		if test.expErrContain == "" {

			if err != nil {
				t.Fatalf("test #%d: error: %v", i, err)
			}

			var res int
			err2 := f.DB.QueryRow(
				"SELECT mark FROM group_event_members WHERE group_eventID =? AND group_membersID = ?",
				test.eventID, test.sel_group_membersID).
				Scan(&res)
			if err2 != nil {
				t.Fatalf("test #%d: error: %v", i, err2)
			}

			if test.exp != res {
				t.Fatalf("test #%d: Expected %v, got %v", i, test.exp, res)
			}
		}

		if test.expErrContain != "" {
			fmt.Printf("Error: %v\n", err)
			if err == nil {
				t.Fatalf("test #%d: Expected error, got nil", i)
			}
			if !strings.Contains(err.Error(), test.expErrContain) {
				t.Fatalf("test #%d: Expected error '%s', got %v", i, test.expErrContain, err.Error())
			}
		}
	}
}
