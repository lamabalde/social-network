package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

/*creates a new chat for the given users */
func (dbm *DBModel) CreatePrivateChat(userID1, userID2 string) (models.Chat, error) {
	var chat models.Chat
	chatID, err := helpers.GenerateNewUUID()
	if err != nil {
		return chat, err
	}

	q := `INSERT INTO chat_members (chatID, userID) VALUES (?,?);
		  INSERT INTO chat_members (chatID, userID) VALUES (?,?);
	`
	err = dbm.runTransations(q, chatID, userID1, chatID, userID2)
	if err != nil {
		return chat, fmt.Errorf("AddUserSToChat run transaction failed: %w", err)
	}

	chat.ID = chatID
	chat.Type = models.CHAT_TYPE_PRIVATE
	return chat, nil
}

/*
inserts a new chat message into DB, returns an ID for the message
*/
func (dbm *DBModel) AddPrivateChatMessage(chatID string, authorID string, content string, images []string, dateCreate time.Time) (string, error) {
	id, err := helpers.GenerateNewUUID()
	if err != nil {
		return "", err
	}

	chatMembersID, err := dbm.getPrivateChatMembersID(chatID, authorID)
	if err == models.ErrNoRecords {
		return "", models.ErrAddConstaints
	}
	if err != nil {
		return "", fmt.Errorf("failed in getting a chatMembersID: %w", err)
	}

	strOfImages := helpers.JoinToNullString(images)

	q := `INSERT INTO chat_messages (id, content, images, chat_membersID, dateCreate) VALUES (?,?,?,?,?)`
	_, err = dbm.DB.Exec(q, id, content, strOfImages, chatMembersID, dateCreate)
	if err != nil {
		return "", fmt.Errorf("AddPrivateChatMessage exec query failed: %w", err)
	}

	return id, nil
}

func (dbm *DBModel) DeletePrivateChatMessage(id string) error {
	q := `DELETE FROM chat_messages WHERE id=?`
	_, err := dbm.DB.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (dbm *DBModel) GetLastMessageDateFromUserToRecipient(userId string, recipientID string) (string, error) {
	q := `SELECT max(ms.dateCreate) FROM chat_messages ms 
	WHERE ms.chat_membersID IN (SELECT  mb.id as mbID FROM chat_members mb
		WHERE mb.userID=? AND mb.chatID IN (SELECT chatID FROM chat_members WHERE userID=?)) 
	GROUP BY ms.chat_membersID `

	var messageDateCreate sql.NullString
	row := dbm.DB.QueryRow(q, recipientID, userId)
	err := row.Scan(&messageDateCreate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messageDateCreate.String, models.ErrNoRecords
		}
		return messageDateCreate.String, err
	}

	return messageDateCreate.String, nil
}

func (dbm *DBModel) GetPrivateChat(usrID1, usrID2 string) (models.Chat, error) {
	var chat models.Chat
	chat.Type = models.CHAT_TYPE_PRIVATE

	q := `SELECT mb1.chatID, users.userName 
			FROM chat_members mb1
			INNER JOIN chat_members mb2 ON mb1.chatID=mb2.chatID AND mb2.userID=?
			LEFT JOIN users ON mb1.userID = users.id
			WHERE mb1.userID=?
	`
	// TODO old query q := `SELECT chatID, users.userName  
			// FROM chat_members mb
			// LEFT JOIN users ON mb.userID = users.id
			// WHERE mb.userID=?
			//   AND chatID IN (SELECT chatID FROM chat_members cmb WHERE cmb.userID=?) 
	// `//usrID2, usrID1

	row := dbm.DB.QueryRow(q, usrID1, usrID2)

	err := row.Scan(&chat.ID, &chat.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return chat, models.ErrNoRecords
		}
		return chat, err
	}

	if chat.ID == "" {
		return chat, models.ErrNoRecords
	}

	return chat, nil
}

func (dbm *DBModel) GetPrivateChatMembersIDs(chatID string) ([]string, error) {
	members := make([]string, 2)

	q := `SELECT userID  
			FROM chat_members mb
			WHERE mb.chatID=?
	`

	rows, err := dbm.DB.Query(q, chatID)
	if err != nil {
		return members, err
	}

	defer rows.Close()

	counter := 0
	for rows.Next() {
		if counter == 2 {
			return members, models.ErrTooManyRecords
		}

		var userID string

		err := rows.Scan(&userID)
		if err != nil {
			return members, err
		}

		members[counter] = userID

		counter++

	}

	if counter < 2 {
		return members, fmt.Errorf("only %d members in private chat", counter)
	}

	return members, nil
}

/*
returns chat with 'limit' messages of the chat with the given 'chatID', skips 'offset' last messages
*/
func (dbm *DBModel) GetPrivateChatMessagesByChatId(chatID string, limit, offset int) ([]*models.ChatMessage, error) {
	condition, arguments := createConditionSelectMessagesByChatId(chatID, limit, offset)
	query := createQuerySelectPrivateChatMessages(
		" chat_messages.id, chat_messages.content, chat_messages.images, chat_messages.dateCreate, chat_members.userID, users.userName ",
		condition,
	)
	chatMessages, err := dbm.execQuerySelectChatMessages(query, arguments, scanRowToChatMessage)
	if err != nil {
		return nil, err
	}

	return chatMessages, nil
}

/*
returns id of the row in chat_members table for a given chat and user
used in InsertChatMessage
*/
func (dbm *DBModel) getPrivateChatMembersID(chatID string, authorID string) (string, error) {
	var id string
	q := `SELECT id FROM chat_members WHERE chatID=? AND userID=? `
	row := dbm.DB.QueryRow(q, chatID, authorID)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNoRecords
		}
		return "", err
	}

	return id, nil
}

/*
creates a condition for the query to select chat messages.
Used in the GetChatMessagesByChatId function.
*/
func createConditionSelectMessagesByChatId(chatID string, limit, offset int) (string, []any) {
	condition := ` WHERE chat_members.chatID = ? `
	arguments := []any{chatID}

	arguments = append(arguments, limit, offset)

	return condition, arguments
}

/*
creates the query to select chat messages.
Used in the GetChatMessagesByChatId and GetChatMessagesByUsersId function.
*/
func createQuerySelectPrivateChatMessages(fields, condition string) (query string) {
	query = `
		SELECT  ` + fields + ` 
	    FROM chat_messages
	    LEFT JOIN chat_members ON chat_messages.chat_membersID=chat_members.id  
		LEFT JOIN users ON chat_members.userID=users.id 
		` + condition + ` 
		ORDER BY chat_messages.dateCreate DESC
		LIMIT ? OFFSET ?;
		`
	return
}

/*
executes a query to select chat messages.
Used in the Get...ChatMessages... functions.
*/
func (dbm *DBModel) execQuerySelectChatMessages(query string, arguments []any, scanRowToChat func(*sql.Rows) (*models.ChatMessage, error)) ([]*models.ChatMessage, error) {
	var chatMessages []*models.ChatMessage

	rows, err := dbm.DB.Query(query, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		message, err := scanRowToChat(rows)
		if err != nil {
			return nil, err
		}
		chatMessages = append(chatMessages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chatMessages, nil
}

/*
scans a row from a query to the the item of models.Chat.
Used in the execQuerySelectMessages function.
*/
func scanRowToChatMessage(rows *sql.Rows) (*models.ChatMessage, error) {
	message := &models.ChatMessage{}
	var images sql.NullString

	// parse the row with fields:
	// SELECT  .id, .content, .images, .dateCreate, _members.userID, users.userName
	err := rows.Scan(&message.ID, &message.Content, &images, &message.DateCreate, &message.UserID, &message.UserName)
	message.Images = helpers.SplitNullString(images)
	return message, err
}
