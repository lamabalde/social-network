package controllers

import (
	"errors"
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wshub"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

/*
sends the list of online users to the current user
and the current user's online status to the other users
*/
func SendOnlineUsers(app *application.Application, currConnection *wsconnection.UsersConnection) error {
	onlineUsers := app.Hub.GetOnlineUsers()

	err := sendOnlineUsersToCurrentUser(app, currConnection, onlineUsers)
	if err != nil && !errors.Is(err, webmodel.ErrWarning) {
		return currConnection.WSError(fmt.Sprintf("sending list of online users to the user %s faild", currConnection.Session.User), err)
	}

	// send the new online user to their followers/ings

	err = SendUserStatusToUsers(app, webmodel.UserOnline)(currConnection, webmodel.WSMessage{})
	if err != nil {
		return errors.Join(webmodel.ErrWarning, err)
	}

	return nil
}

/*
sends the list of users in chat to the new joined user
and the new user's joined status to the other users in chat
*/
func SendChattingUsers(app *application.Application, currConnection *wsconnection.UsersConnection) error {
	onlineUsers := app.Hub.GetOnlineUsers()

	err := sendUsersInChatToCurrentUser(app, currConnection, onlineUsers)
	if err != nil && !errors.Is(err, webmodel.ErrWarning) {
		return currConnection.WSError(fmt.Sprintf("sending list of users in chat to the user %s faild", currConnection.Session.User), err)
	}

	// send the new user to the chat members

	err = SendUserStatusToChatMembers(app, webmodel.UserJoinedChat)(currConnection, webmodel.WSMessage{})
	if err != nil {
		return errors.Join(webmodel.ErrWarning, err)
	}

	return nil
}

/*
sends list of online users
*/
func sendOnlineUsersToCurrentUser(app *application.Application, currConnection *wsconnection.UsersConnection, onlineUsers wshub.MapID) error {
	users, err := app.DBModel.GetFilteredFollowsOrderedByMessagesToGivenUser(onlineUsers, currConnection.Client.UserID)
	if err != nil {
		return currConnection.WSError("get the users from DB failed", err)
	}

	_, err = currConnection.SendSuccessMessage(webmodel.OnlineUsers, users)
	return err
}

/*
sends list of users in chat to the current user
*/
func sendUsersInChatToCurrentUser(app *application.Application, currConnection *wsconnection.UsersConnection, onlineUsers wshub.MapID) error {
	users := currConnection.Client.Room.GetUsersInRoom()

	_, err := currConnection.SendSuccessMessage(webmodel.ChattingUsers, users)
	return err
}

/*
sends current user's status (on/offline) to their follows
*/
func SendUserStatusToUsers(app *application.Application, statusType string) wsconnection.FuncReplier {
	return func(currConnection *wsconnection.UsersConnection, wsMessage webmodel.WSMessage) error {
		currentUser := models.UserBase{ID: currConnection.Client.UserID, UserName: currConnection.Client.UserName}

		users, err := app.DBModel.GetFollows(currConnection.Client.UserID)
		if err != nil {
			return currConnection.WSError("getFollows from DB failed", err)
		}

		usersID := make([]string, len(users))
		for i, user := range users {
			usersID[i] = user.ID
		}

		_, _, err = currConnection.SendMessageToUsers(usersID, statusType, currentUser)
		if err != nil {
			return currConnection.WSError("SendMessageToUsers failed", err)
		}
		return nil
	}
}

/*
sends current user's status (join/quit chat) to users in the user's chat room
*/
func SendUserStatusToChatMembers(app *application.Application, statusType string) wsconnection.FuncReplier {
	return func(currConnection *wsconnection.UsersConnection, wsMessage webmodel.WSMessage) error {
		currentUser := models.UserBase{ID: currConnection.Client.UserID, UserName: currConnection.Client.UserName}

		_, _, err := currConnection.SendMessageToClientRoom(statusType, currentUser)
		if err != nil {
			return currConnection.WSError("SendMessageToClientRoom failed", err)
		}
		return nil
	}
}
