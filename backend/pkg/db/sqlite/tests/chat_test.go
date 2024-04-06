package tests

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
)

func TestGetChatMessagesByUsersIds(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var chat models.Chat

	f := queries.DBModel{DB: db}
	fmt.Println("--------------chat between 2 & 3 (no chat)------------------")
	chat, err = f.GetPrivateChat("2", "3")
	if err != models.ErrNoRecords {
		t.Fatalf("expected %v, got %v", models.ErrNoRecords, err)
	}

	fmt.Println("chat=", chat)

	fmt.Println("---------chat between 1 & 2------------------")
	chat, err = f.GetPrivateChat("1", "2")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("chat=", chat)
	expRes := models.Chat{
		ID:   "chat1",
		Name: "test1",
		Type: 0,
	}
	if !reflect.DeepEqual(expRes, chat) {
		t.Fatalf("expected \n%#v\n got \n%#v", expRes, chat)
	}

	chat.Messages, err = f.GetPrivateChatMessagesByChatId(chat.ID, 3, 0)
	if err != nil {
		t.Fatal(err)
	}
	expmess := []*models.ChatMessage{
		{
			ID:         "7",
			Content:    "mes3 from 2 to 1",
			UserID:     "2",
			UserName:   "test1",
			DateCreate: time.Date(2023, time.November, 22, 10, 59, 33, 656479916, time.UTC),
		},
		{
			ID:         "6",
			Content:    "mes4 from 1 to 2",
			UserID:     "1",
			UserName:   "no",
			DateCreate: time.Date(2023, time.November, 22, 10, 58, 33, 656479916, time.UTC),
		},
		{
			ID:         "5",
			Content:    "mes2 from 2 to 1",
			UserID:     "2",
			UserName:   "test1",
			DateCreate: time.Date(2023, time.November, 22, 10, 58, 23, 656479916, time.UTC),
		},
	}

	expRes.Messages=  expmess
	
	if !reflect.DeepEqual(expRes, chat) {
		t.Fatalf("Expected\n %v, got\n %v", &expRes, &chat)
	}

	fmt.Printf("%v\n", chat)

	offset := 4
	fmt.Println("-----get 1 message offset -", offset)
	chat.Messages, err = f.GetPrivateChatMessagesByChatId(chat.ID, 1, offset)
	if err != nil {
		t.Fatal(err)
	}

	expmess = []*models.ChatMessage{
		{
			ID:         "3",
			Content:    "mes2 from 1 to 2",
			UserID:     "1",
			UserName:   "no",
			DateCreate: time.Date(2023, time.November, 22, 10, 56, 23, 656479916, time.UTC),
		},
	}

	expRes.Messages=  expmess
	if !reflect.DeepEqual(expRes, chat) {
		t.Fatalf("Expected\n %v, got\n %v", expRes, chat)
	}
	fmt.Printf("%v\n", chat)

	offset = 7
	fmt.Println("-----get 2 messages offset (0 messages) - ", offset)
	chat.Messages, err = f.GetPrivateChatMessagesByChatId(chat.ID, 2, offset)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", chat)

	fmt.Println("-----------------chat between 1 & 3 (no messages )----------------")
	chat, err = f.GetPrivateChat("1", "3")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("id=", chat)
	expRes = models.Chat{
		ID:   "chat2",
		Name: "test2",
		Type: 0,
	}
	if !reflect.DeepEqual(expRes, chat) {
		t.Fatalf("expected %v, got %v", expRes, chat)
	}
	fmt.Println("-----get last 3 messages")
	expmess = []*models.ChatMessage{
		{
			ID:         "9",
			Content:    "mes1 from 1 to 3",
			UserID:     "1",
			UserName:   "no",
			DateCreate: time.Date(2023, time.November, 24, 10, 59, 33, 656479916, time.UTC),
		},
	
		{
			ID:         "8",
			Content:    "hello from 1 to 3",
			UserID:     "1",
			UserName:   "no",
			DateCreate: time.Date(2023, time.November, 23, 10, 59, 33, 656479916, time.UTC),
		},
	}
	expRes.Messages=expmess
	chat.Messages, err = f.GetPrivateChatMessagesByChatId(chat.ID, 3, 0)
	if err != nil {
		t.Fatal(err)
	}


	if !reflect.DeepEqual(expRes, chat) {
		t.Fatalf("Expected\n %v,\n got\n %v", expRes, chat)
	}
	fmt.Printf("%v\n", chat)
}

func TestGetLastMessageDateFromUserToRecipient(t *testing.T) {
	db, err := sqlite.OpenDatabase(T_DBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	f := queries.DBModel{DB: db}

	fmt.Println("-------get last date of message from user 1 to user 2 -----------")
	date, err := f.GetLastMessageDateFromUserToRecipient("1", "2")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	fmt.Println("---------------")

	expRes := "2023-11-22 10:59:33.656479916+00:00"
	if date != expRes {
		t.Fatalf("expected date=%v, got %v", expRes, date)
	}
	fmt.Printf("from user %d to user %d mes date: %v\n", 2, 3, date)

	fmt.Println("-------get last date of message from user  to user 1 -----------")
	date, err = f.GetLastMessageDateFromUserToRecipient("2", "1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	fmt.Println("---------------")

	expRes = "2023-11-22 10:58:33.656479916+00:00"
	if date != expRes {
		t.Fatalf("expected date=%v, got %v", expRes, date)
	}
	fmt.Printf("from user %d to user %d mes date: %v\n", 2, 3, date)
}
