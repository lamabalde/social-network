package tests

import (
	"database/sql"
	"testing"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

func TestAddDeleteNotification(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	numofRows, err := getRowNumbersInTable(f.DB, "notifications", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add ----------------------------------------------------------------
	tests := []struct {
		note    models.Notification
		exprows int
	}{
		{
			note: models.Notification{
				UserID:     "2",
				FromUserID: sql.NullString{"3", true},
				Type:       "message",
				Body:       "new message",
				Read:       1,
				DateCreate: time.Date(2023, time.March, 2, 11, 11, 11, 0, time.UTC),
			},
			exprows: numofRows + 1,
		},

		{
			note: models.Notification{
				UserID:     "3",
				PostID:     sql.NullString{"2", true},
				Type:       "post comment",
				Body:       "new comment",
				Read:       1,
				DateCreate: time.Date(2023, time.March, 3, 11, 11, 11, 0, time.UTC),
			},
			exprows: numofRows + 2,
		},
	}
	ids := make([]int, len(tests))
	for i, test := range tests {
		id, err := f.AddNotification(&test.note)
		if err != nil {
			t.Fatalf("test #%d: error: %v", i, err)
		}
		numofRowsres, err2 := getRowNumbersInTable(f.DB, "notifications", "", nil)
		if err2 != nil {
			t.Fatal(err)
		}

		if test.exprows != numofRowsres {
			t.Fatalf("\ntest #%d: Expected %v, got %v", i, test.exprows, numofRowsres)
		}

		ids[i] = id
	}

	for i := 1; i < len(ids); i++ {
		if ids[i] != ids[i-1]+1 {
			t.Fatalf("wrong id for #%d: Expected %d, got %d", i, ids[i-1]+1, ids[i])
		}
	}

	// Delete ---------------------------------------------------------------

	for i := 0; i < len(ids); i++ {
		err = f.DeleteNotification(ids[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	numofRowsres, err := getRowNumbersInTable(f.DB, "notifications", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if numofRowsres != numofRows {
		t.Fatalf("Expected %v, got %v", numofRows, numofRowsres)
	}
}

func TestGetUserUnReadNotification(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	tests := []struct {
		userID string
		expRes []*models.Notification
	}{
		{
			userID: "1",
			expRes: []*models.Notification{
				{
					ID:         1,
					UserID:     "1",
					FromUserID: sql.NullString{"3", true},
					Type:       "",
					Body:       "text",
					Read:       0,
					DateCreate: time.Date(2023, time.November, 25, 10, 55, 23, 656479916, time.UTC),
				},
				{
					ID:         3,
					UserID:     "1",
					FromUserID: sql.NullString{"3", true},
					Type:       "follow request",
					Body:       "test2 wants to follow you",
					Read:       0,
					DateCreate: time.Date(2023, time.December, 12, 10, 55, 23, 656479916, time.UTC),
				},
			},
		},

		{
			userID: "2",
			expRes: nil,
		},
	}
	for i, test := range tests {
		notes, err := f.GetUserUnReadNotification(test.userID)
		if err != nil {
			t.Fatalf("test #%d: error: %v", i, err)
		}
		if test.expRes == nil {
			if notes != nil {
				t.Fatalf("\ntest #%d: Expected %v, got %v", i, test.expRes, notes)
			}
			continue
		}
		if len(test.expRes) != len(notes) {
			t.Fatalf("\ntest #%d: Expected %d, got %d", i, len(test.expRes), len(notes))
		}
		for j := 0; j < len(test.expRes); j++ {
			if test.expRes[j].ID != notes[j].ID ||
				test.expRes[j].UserID != notes[j].UserID ||
				test.expRes[j].FromUserID != notes[j].FromUserID ||
				test.expRes[j].Type != notes[j].Type ||
				test.expRes[j].Body != notes[j].Body ||
				test.expRes[j].Read != notes[j].Read {
				t.Fatalf("\ntest #%d elm#%d: Expected %v, got %v", i, j, test.expRes[j], notes[j])
			}
		}
	}
}
