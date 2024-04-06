package tests

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

func TestCreateGroup(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	loc1, err := time.LoadLocation("Europe/Kyiv")
	if err != nil {
		t.Fatal(err)
	}

	loc2 := time.FixedZone("UTC-0800", -8*60*60)

	newGroups := []models.Group{
		{
			Name:        "Group Title",
			Description: "group description",
			CreatorID:   "1",
			DateCreate:  time.Date(2023, time.March, 1, 11, 11, 11, 0, loc1),
		},
		{
			Name:        "Group2 Title",
			Description: "group2 description",
			CreatorID:   "2",
			DateCreate:  time.Date(2023, time.March, 2, 12, 12, 12, 0, loc2),
		},
	}
	ids := make([]string, len(newGroups))
	for i, group := range newGroups {
		ids[i], err = f.CreateGroup(group.Name, group.Description, group.CreatorID, group.DateCreate)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, id := range ids {
		rows, err := f.DB.Query(`SELECT * FROM groups WHERE id =?`, id)
		if err != nil {
			t.Fatal(err)
		}
		selected := &models.Group{}
		counter := 0
		for rows.Next() {
			if counter == 1 {
				t.Fatal("too many rows")
			}

			err := rows.Scan(&selected.ID, &selected.Name, &selected.Description, &selected.CreatorID, &selected.DateCreate)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("group date created: ", newGroups[i].DateCreate)
			fmt.Println("selected date created: ", selected.DateCreate)
			counter++
		}

		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		rows.Close()

		if id != selected.ID ||
			newGroups[i].Name != selected.Name ||
			newGroups[i].Description != selected.Description ||
			newGroups[i].CreatorID != selected.CreatorID || newGroups[i].DateCreate.Local() != selected.DateCreate.Local() {
			t.Fatalf("Expected %s with id %s,\ngot %s", &newGroups[i], id, selected)
		}
	}
}

func TestAddDelUserToGroup(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	q := `select count(*) from group_members`
	rowsNumberBefore := 0
	row := f.DB.QueryRow(q)
	row.Scan(&rowsNumberBefore)

	new := struct{ groupID, userID string }{"3", "2"}
	err = f.AddUserToGroup(new.groupID, new.userID, true)
	if err != nil {
		t.Fatal(err)
	}

	rowsNumberAfter := 0
	row = f.DB.QueryRow(q)
	row.Scan(&rowsNumberAfter)
	if rowsNumberAfter-rowsNumberBefore != 1 {
		t.Fatal(err)
	}

	err = f.DeletUserFromGroup(new.groupID, new.userID)
	if err != nil {
		t.Fatal(err)
	}

	row = f.DB.QueryRow(q)
	row.Scan(&rowsNumberAfter)
	if rowsNumberAfter != rowsNumberBefore {
		t.Fatal(err)
	}
}

func TestAddDelUserMessage(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	loc1, err := time.LoadLocation("Europe/Kyiv")
	if err != nil {
		t.Fatal(err)
	}
	type expT struct {
		groupID, authorID, text string
		group_membersID         int64
		images                  []string
		dateCreate              time.Time
	}

	exp := expT{
		groupID:         "1",
		authorID:        "2",
		text:            "message text",
		group_membersID: 1,
		images:          nil,
		dateCreate:      time.Date(2023, time.March, 1, 11, 11, 11, 0, loc1),
	}
	id, err := f.AddGroupChatMessage(exp.groupID, exp.authorID, exp.text, exp.images, exp.dateCreate)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := f.DB.Query(`SELECT * FROM group_chat WHERE id =?`, id)
	if err != nil {
		t.Fatal(err)
	}
	selected := struct {
		id              string
		content         string
		images          sql.NullString
		group_membersID int64
		dateCreate      time.Time
	}{}
	counter := 0
	for rows.Next() {
		if counter == 1 {
			t.Fatal("too many rows")
		}

		err := rows.Scan(&selected.id, &selected.content, &selected.images, &selected.group_membersID, &selected.dateCreate)
		if err != nil {
			t.Fatal(err)
		}
		counter++
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	rows.Close()

	if id != selected.id ||
		exp.text != selected.content ||
		selected.images.Valid ||
		exp.group_membersID != selected.group_membersID ||
		exp.dateCreate.Local() != selected.dateCreate.Local() {
		t.Fatalf("Expected %#v with id %s,\ngot %#v", &exp, id, selected)
	}
	fmt.Println("selected.group_membersID: ", selected.group_membersID)

	// ----- delete the message -------
	err = f.DeleteGroupChatMessage(id)
	if err != nil {
		t.Fatal(err)
	}

	err = f.DB.QueryRow(`SELECT id FROM group_chat WHERE id =?`, id).Scan()
	if err != sql.ErrNoRows {
		t.Fatal(err)
	}

	// ----- must throw an error (the user is not in the group)---------

	exp = expT{
		groupID:    "3",
		authorID:   "3",
		text:       "message text",
		images:     nil,
		dateCreate: time.Date(2023, time.March, 1, 11, 11, 11, 0, loc1),
	}
	_, err = f.AddGroupChatMessage(exp.groupID, exp.authorID, exp.text, exp.images, exp.dateCreate)
	if err == nil {
		t.Fatal("must throw an error")
	}
	fmt.Println("Error is:", err)
}

func TestGetListOfMembers(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	groups := []models.Group{
		{
			ID: "2",
			Members: []models.UserBase{
				{ID: "1", UserName: "no"},
				{ID: "2", UserName: "test1"},
				//{ID: "3", UserName: "test2"},
			},
		},
		{
			ID: "1",
			Members: []models.UserBase{
				{ID: "2", UserName: "test1"},
				{ID: "3", UserName: "test2"},
			},
		},
	}

	for _, group := range groups {
		members, err := f.GetListOfGroupMembers(group.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(members, group.Members) {
			t.Fatalf("Expected %v, got %v", group.Members, members)
		}
	}
}

func TestGetListOfUserGroups(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	testCases := []struct {
		userID string
		groups []models.Group
	}{
		{
			userID: "1",
			groups: []models.Group{
				{ID: "2", Name: "group2"},
				{ID: "3", Name: "group3"},
			},
		},
		{
			userID: "2",
			groups: []models.Group{
				{ID: "1", Name: "group1"},
				{ID: "2", Name: "group2"},
			},
		},
		{
			userID: "3",
			groups: []models.Group{
				{ID: "1", Name: "group1"},
				//{ID: "2", Name: "group2"},
			},
		},
	}

	for i, test := range testCases {
		groups, err := f.GetListOfUserGroups(test.userID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(groups, test.groups) {
			t.Fatalf("test# %d: Expected %v, got %v", i, test.groups, groups)
		}
	}
}

func TestGetGroupChatMessagesByGroupId(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}
	testCases := []struct {
		groupID       string
		limit, offset int
		messages      []*models.ChatMessage
	}{
		{
			groupID: "1",
			limit:   10,
			offset:  0,
			messages: []*models.ChatMessage{
				{ID: "7", UserID: "3", UserName: "test2", Content: "mess1 from 3 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 59, 33, 656479916, time.UTC)},
				{ID: "4", UserID: "2", UserName: "test1", Content: "mess3 from 2 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 35, 23, 656479916, time.UTC)},
				{ID: "3", UserID: "2", UserName: "test1", Content: "mess2 from 2 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 25, 23, 656479916, time.UTC)},
				{ID: "2", UserID: "3", UserName: "test2", Content: "hello from 3 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 15, 23, 656479916, time.UTC)},
				{ID: "1", UserID: "2", UserName: "test1", Content: "hello from 2 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 0o5, 23, 656479916, time.UTC)},
			},
		},
		{
			groupID: "1",
			limit:   10,
			offset:  3,
			messages: []*models.ChatMessage{
				{ID: "2", UserID: "3", UserName: "test2", Content: "hello from 3 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 15, 23, 656479916, time.UTC)},
				{ID: "1", UserID: "2", UserName: "test1", Content: "hello from 2 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 0o5, 23, 656479916, time.UTC)},
			},
		},
		{
			groupID: "1",
			limit:   2,
			offset:  0,
			messages: []*models.ChatMessage{
				{ID: "7", UserID: "3", UserName: "test2", Content: "mess1 from 3 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 59, 33, 656479916, time.UTC)},
				{ID: "4", UserID: "2", UserName: "test1", Content: "mess3 from 2 in 1", DateCreate: time.Date(2023, time.November, 22, 10, 35, 23, 656479916, time.UTC)},
			},
		},
		{
			groupID: "2",
			limit:   10,
			offset:  0,
			messages: []*models.ChatMessage{
				{ID: "6", UserID: "1", UserName: "no", Content: "mess1 from 1 in 2", DateCreate: time.Date(2023, time.November, 22, 10, 55, 33, 656479916, time.UTC)},
				//{ID: "5", UserID: "2", UserName: "test1", Content: "mess1 from 2 in 2", DateCreate: time.Date(2023, time.November, 22, 10, 45, 23, 656479916, time.UTC)},
			},
		},
		{
			groupID:  "3",
			limit:    10,
			offset:   0,
			messages: nil,
		},
	}

	for i, test := range testCases {
		messages, err := f.GetGroupChatMessagesByGroupId(test.groupID, test.limit, test.offset)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(test.messages, messages) {
			t.Fatalf("test#%d: Expected %v, got %v", i, test.messages, messages)
		}
		if len(test.messages) != len(messages) {
			t.Fatalf("test#%d: Expected %d messages, got  %d messages", i, len(test.messages), len(messages))
		}
		for j, message := range test.messages {
			if !reflect.DeepEqual(test.messages[j], message) {
				t.Fatalf("test#%d: Expected %v, got %v", i, test.messages[j], message)
			}
		}
	}
}

func TestGetGroupByID(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	testCases := []struct {
		groupID string
		group   models.Group
	}{
		{
			groupID: "1",
			group: models.Group{
				ID: "1", Name: "group1", Description: "group1 by user 2, id 1", CreatorID: "2",
				DateCreate: time.Date(2023, time.November, 23, 10, 55, 23, 656479916, time.UTC),
				Members: []models.UserBase{
					{ID: "2", UserName: "test1"},
					{ID: "3", UserName: "test2"},
				},
			},
		},
		{
			groupID: "2",
			group: models.Group{
				ID: "2", Name: "group2", Description: "group2 by user 2, id 2", CreatorID: "2",
				DateCreate: time.Date(2023, time.November, 24, 10, 55, 23, 656479916, time.UTC),
				Members: []models.UserBase{
					{ID: "1", UserName: "no"},
					// {ID: "2", UserName: "test1"},
					// {ID: "3", UserName: "test2"},
				},
			},
		},
		{
			groupID: "3",
			group: models.Group{
				ID: "3", Name: "group3", Description: "group3 by user 1, id 3", CreatorID: "1",
				DateCreate: time.Date(2023, time.November, 25, 10, 55, 23, 656479916, time.UTC),
				Members: []models.UserBase{
					{ID: "1", UserName: "no"},
				},
			},
		},
	}

	for i, test := range testCases {
		group, err := f.GetGroupByID(test.groupID)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(group, test.group) {
			t.Fatalf("test#%d, get group: Expected %v, got %v", i, &test.group, &group)
		}

		groupName, err := f.GetGroupNameByID(test.groupID)
		if err != nil {
			t.Fatal(err)
		}
		if groupName != test.group.Name {
			t.Fatalf("test#%d, get group name: Expected %v, got %v", i, test.group.Name, groupName)
		}
	}
}

/*
func TestGetGroupMessagesByUsersIds(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var group *models.Group

	f := queries.DBModel{DB: db}

	fmt.Println("------------get group 3 members--------------------")
	groupChat, name, err := f.GetPrivateGroup("2", "3")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("------------get group 3 chat (no messages) --------------------")

	fmt.Println("groupChat=", groupChat)
	group, err = f.GetPrivateGroupMessagesByGroupId(groupChat, 3, 0)
	if err != nil {
		t.Fatal(err)
	}

	group.ID = groupChat
	group.Name = name
	fmt.Printf("%s\n", group.String())

	fmt.Println("------------get group 1 members--------------------")
	groupChat, name, err = f.GetPrivateGroup("4", "5")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("groupChat=", groupChat)
	fmt.Println("-----get last 3 messages")
	if groupChat != "" {
		group, err = f.GetPrivateGroupMessagesByGroupId(groupChat, 3, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	group.ID = groupChat
	group.Name = name
	fmt.Printf("%s\n", group.String())

	offset := 2
	fmt.Println("-----get 3 messages offset ", offset)
	if groupChat != "" {
		group, err = f.GetPrivateGroupMessagesByGroupId(groupChat, 3, offset)
		if err != nil {
			t.Fatal(err)
		}
	}
	group.ID = groupChat
	group.Name = name
	fmt.Printf("%s\n", group.String())

	offset = 3
	fmt.Println("-----get 3 messages offset ", offset)
	if groupChat != "" {
		group, err = f.GetPrivateGroupMessagesByGroupId(groupChat, 3, offset)
		if err != nil {
			t.Fatal(err)
		}
	}
	group.ID = groupChat
	group.Name = name
	fmt.Printf("%s\n", group.String())
}
*/
