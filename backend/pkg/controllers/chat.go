package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

const (
	PRIVATE_CHAT_ROOM = "private"
	GROUP_CHAT_ROOM   = "group"
)

/*
sends list of the current user's followers and followings and joined groups

// TODO: sort chat list based on last message
*/
func ReplyGetChatList(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		follows, err := app.DBModel.GetFollows(currConnection.Client.UserID)
		if err != nil {
			return nil, currConnection.WSError("ReplyGetChatList: Failed to get search data from database", err)
		}
		groups, _ := app.DBModel.GetListOfUserGroups(currConnection.Client.UserID)
		chats := []models.Chat{}
		for _, user := range follows { // list of followers/following
			chat := models.Chat{ID: user.ID, Name: user.UserName, Type: models.CHAT_TYPE_PRIVATE}
			chats = append(chats, chat)
		}
		for _, group := range groups { // list of user's groups
			chat := models.Chat{ID: group.ID, Name: group.Name, Type: models.CHAT_TYPE_GROUP}
			chats = append(chats, chat)
		}
		return chats, nil
	}
}

/*
opens a new chat between user and recepient.
used in OpenChat handler (to request /openChat/private?id=<recepientID>, userID will be the current user's ID)
*/
func OpenPrivateChat(app *application.Application, userID, recepientID string) (models.Chat, error) {
	chat, err := app.DBModel.GetPrivateChat(userID, recepientID)
	if errors.Is(err, models.ErrNoRecords) {
		return newPrivateChat(app, userID, recepientID)
	}
	if err != nil {
		return models.Chat{}, fmt.Errorf("OpenPrivateChat: get the chat id from DB failed: %v", err)
	}

	chat.Type = models.CHAT_TYPE_PRIVATE

	chat.Messages, err = app.DBModel.GetPrivateChatMessagesByChatId(chat.ID, webmodel.CHAT_MESSAGES_PORTION, 0)
	if err != nil {
		return chat, fmt.Errorf("OpenPrivateChat: get the chat messages from DB failed: %v", err)
	}
	slices.Reverse(chat.Messages)

	return chat, nil
}

/*
	gets group chat from DB along with its CHAT_MESSAGES_PORTION(10) last messages.

If there is no group for the given groupID, it returns not nil error
*/
func OpenGroupChat(app *application.Application, groupID string) (models.Chat, error) {
	groupChat := models.Chat{ID: groupID, Type: models.CHAT_TYPE_GROUP}
	var err error
	groupChat.Name, err = app.DBModel.GetGroupNameByID(groupID)
	if errors.Is(err, models.ErrNoRecords) {
		return models.Chat{}, fmt.Errorf("OpenGroupChat: no group with id %s", groupID)
	}
	if err != nil {
		return models.Chat{}, fmt.Errorf("OpenGroupChat: get the group id from DB failed: %v", err)
	}

	groupChat.Messages, err = app.DBModel.GetGroupChatMessagesByGroupId(groupID, webmodel.CHAT_MESSAGES_PORTION, 0)
	if err != nil {
		return groupChat, fmt.Errorf("OpenGroupChat: get the group chat messages from DB failed: %v", err)
	}
	slices.Reverse(groupChat.Messages)

	return groupChat, nil
}

/*
broadcasts messages to clients connected to the chat
*/
func ReplySendMessageToChat(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		chatMessage, err := parse.PayloadToChatMessage(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a chat message: '%s'", message.Payload), err)
		}

		errmessage := chatMessage.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		chatMembers, err := getChatMembersFromDB(app, currConnection)
		if err != nil {
			return nil, currConnection.WSError("get chat members from DB failed", err)
		}

		messageId, err := addMessageInDB(app, currConnection, chatMessage)
		if err != nil {
			return nil, currConnection.WSError("get chat members from DB failed", err)
		}

		chatMessage = fillInChatMessage(app, currConnection, chatMessage, messageId, chatMembers) // create a new ws message with type webmodel.InputChatMessage

		jsonMessage, usersSentTo, err := currConnection.SendMessageToClientRoom(webmodel.InputChatMessage, chatMessage)
		if err != nil {
			return nil, currConnection.WSError("sending message to client room failed", err)
		}
		err = sentMessageToUsersNotInChat(app, currConnection, message, chatMembers, usersSentTo, jsonMessage)
		return "sent", err // reply to current user
	}
}

func sentMessageToUsersNotInChat(app *application.Application, currConnection *wsconnection.UsersConnection, message webmodel.WSMessage,
	chatMembers []string, usersSentTo map[string]bool, jsonMessage json.RawMessage,
) error {
	var memebersNotSentTo []string
	for _, memberID := range chatMembers {
		if memberID == currConnection.Client.UserID {
			continue
		}
		if _, ok := usersSentTo[memberID]; !ok {
			memebersNotSentTo = append(memebersNotSentTo, memberID)
		}

	}
	if len(memebersNotSentTo) > 0 {
		// the same message is sent to the "universal" WS connection for users	who don't open chat,
		// so type webmodel.InputChatMessage must be handled in both type of wsconnection (chat and "universal")
		usersSentTo = currConnection.SendBytesToUsers(memebersNotSentTo, jsonMessage)

		var memebersNotSentTo2 []string
		for _, memberID := range memebersNotSentTo {
			if _, ok := usersSentTo[memberID]; !ok {
				memebersNotSentTo2 = append(memebersNotSentTo2, memberID)
			}
		}

		for _, userID := range memebersNotSentTo2 {
			_, err := createNotificationFromUser(app, userID, currConnection.Client.UserID,
				webmodel.NOTE_NEW_PRIVATE_MESSAGE, fmt.Sprintf("you have a new message in the chat with %s", currConnection.Client.UserName))
			if err != nil {
				return handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", userID), err)
			}
		}
	}
	return nil
}

func getChatMembersFromDB(app *application.Application, currConnection *wsconnection.UsersConnection) ([]string, error) {
	currentRoom := currConnection.Client.Room
	var chatMembers []string
	var err error
	switch currentRoom.Type {
	case PRIVATE_CHAT_ROOM:
		chatMembers, err = app.DBModel.GetPrivateChatMembersIDs(currentRoom.ID)
		if err != nil {
			return nil, err
		}

	case GROUP_CHAT_ROOM:
		chatMembers, err = app.DBModel.GetGroupMembersIDs(currentRoom.ID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("wrong type of chat room")
	}
	return chatMembers, nil
}

func addMessageInDB(app *application.Application, currConnection *wsconnection.UsersConnection, chatMessage webmodel.ChatMessage) (string, error) {
	currentRoom := currConnection.Client.Room
	var messageId string
	var err error
	switch currentRoom.Type {
	case PRIVATE_CHAT_ROOM:
		messageId, err = app.DBModel.AddPrivateChatMessage(currentRoom.ID, currConnection.Client.UserID, chatMessage.Content, nil, chatMessage.DateCreate)
		if err != nil {
			return "", err
		}

	case GROUP_CHAT_ROOM:
		messageId, err = app.DBModel.AddGroupChatMessage(currentRoom.ID, currConnection.Client.UserID, chatMessage.Content, nil, chatMessage.DateCreate)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("wrong type of chat room")
	}
	return messageId, nil
}

func ReplyChatPortion(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		offset, err := parse.PayloadToInt(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for the portion of messages '%s'", message.Payload), err)
		}

		currentRoom := currConnection.Client.Room

		chat := models.Chat{
			ID: currentRoom.ID,
		}

		switch currentRoom.Type {
		case PRIVATE_CHAT_ROOM:
			chat.Type = models.CHAT_TYPE_PRIVATE

			chat.Messages, err = app.DBModel.GetPrivateChatMessagesByChatId(currConnection.Client.Room.ID, webmodel.CHAT_MESSAGES_PORTION, offset)
			if err != nil {
				return nil, currConnection.WSError("get the next portion of chat messages from DB failed", err)
			}

		case GROUP_CHAT_ROOM:
			chat.Type = models.CHAT_TYPE_GROUP

			chat.Messages, err = app.DBModel.GetGroupChatMessagesByGroupId(currConnection.Client.Room.ID, webmodel.CHAT_MESSAGES_PORTION, offset)
			if err != nil {
				return nil, currConnection.WSError("get the next portion of chat messages from DB failed", err)
			}

		default:
			return nil, errors.New("wrong type of chat room")
		}

		slices.Reverse(chat.Messages)

		app.InfoLog.Printf("Send the next portion of messages to private chat id: '%s'.", currConnection.Client.Room.ID)

		return chat, nil
	}
}

/*
creates and saves in DB a new chat.
Used in OpenPrivateChat function
*/
func newPrivateChat(app *application.Application, userID, recepientID string) (models.Chat, error) {
	chat, err := app.DBModel.CreatePrivateChat(userID, recepientID)
	if err != nil {
		return chat, fmt.Errorf("creating a chat between users id '%s' and '%s' failed: %v", userID, recepientID, err)
	}

	app.InfoLog.Printf("chat between users id '%s' and '%s' created with id %s", userID, recepientID, chat.ID)
	return chat, nil
}

/*
sends a chat message to the users in chat
used in ReplySendMessageToChat
*/
func fillInChatMessage(app *application.Application, currConnection *wsconnection.UsersConnection, chatMessage webmodel.ChatMessage, messageID string, chatMembers []string) webmodel.ChatMessage {
	chatMessage.MessageID = messageID
	chatMessage.UserID = currConnection.Client.UserID
	chatMessage.UserName = currConnection.Client.UserName
	var chatGenericID string
	if currConnection.Client.Room.Type == GROUP_CHAT_ROOM {
		chatGenericID = currConnection.Client.Room.ID
	} else if currConnection.Client.Room.Type == PRIVATE_CHAT_ROOM {
		chatGenericID = chatMessage.UserID
	}
	chatMessage.GenericID = chatGenericID
	return chatMessage
}
