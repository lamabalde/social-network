package controllers

import (
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyRequestToJoinToGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		groupID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for groupID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(groupID) {
			return nil, currConnection.WSBadRequest(message, "groupID is empty")
		}

		group, err := app.DBModel.GetGroupByID(groupID)
		if err != nil {
			return nil, handleErrGetFromDBByID(app, currConnection, message, "group", groupID, err)
		}

		err = app.DBModel.AddUserToGroup(groupID, currConnection.Client.UserID, false)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message,
				fmt.Sprintf("user (id '%s') to group (id '%s')", currConnection.Client.UserID, groupID), err)
		}

		note, err := createNotificationForGroup(app,
			group.CreatorID,
			currConnection.Client.UserID,
			webmodel.NOTE_JOIN_GROUP_REQUEST,
			fmt.Sprintf("user  %s wants to join group %s", currConnection.Client.UserName, group.Name),
			groupID,
		)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", group.CreatorID), err)
		}

		_, _, err = currConnection.SendMessageToUser(group.CreatorID, webmodel.NewNotification, note)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("sending notification to user (id'%s') failed", group.CreatorID), err)
		}

		return groupID, nil
	}
}

func ReplyInviteToGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userGroup, err := parse.PayloadToGroupUser(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		errmessage := userGroup.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}
		group, err := app.DBModel.GetGroupByID(userGroup.GroupID)
		if err != nil {
			return nil, handleErrGetFromDBByID(app, currConnection, message, "group", userGroup.GroupID, err)
		}
		err = app.DBModel.AddUserToGroup(group.ID, userGroup.User, false)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message,
				fmt.Sprintf("user (id '%s') to group (id '%s')", userGroup.User, group.ID), err)
		}
		note, err := createNotificationForGroup(app,
			userGroup.User,
			currConnection.Client.UserID,
			webmodel.NOTE_INVITE_TO_GROUP,
			fmt.Sprintf("user %s invites you to join group %s", currConnection.Client.UserName, group.Name),
			group.ID,
		)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", userGroup), err)
		}
		_, _, err = currConnection.SendMessageToUser(userGroup.User, webmodel.NewNotification, note)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("sending notification to user (id'%s') failed", userGroup.User), err)
		}
		return userGroup.User, nil
	}
}

func ReplyAcceptRequestToJoinGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		note, err := handleJoinGroupNotification(app, currConnection, message, webmodel.NOTE_JOIN_GROUP_REQUEST)
		if err != nil {
			return nil, err // handleJoinGroupNotification sends errors to wsconnection
		}

		err = app.DBModel.SetGroupMemberStatus(note.GroupID.String, note.FromUserID.String, true)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("set true member status for user (id '%s') in group (id '%s') failed", note.FromUserID.String, note.GroupID.String), err)
		}

		newNote, err := createNotificationForGroup(app,
			note.FromUserID.String,
			currConnection.Client.UserID,
			webmodel.NOTE_JOIN_GROUP_REQUEST_ACCEPTED,
			"your request to join a group has been accepted",
			note.GroupID.String,
		)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", note.FromUserID.String), err)
		}

		_, _, err = currConnection.SendMessageToUser(note.FromUserID.String, webmodel.NewNotification, newNote)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("sending notification to user (id'%s') failed", note.FromUserID.String), err)
		}

		return note.GroupID.String, nil
	}
}

func ReplyDeclineRequestToJoinGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		note, err := handleJoinGroupNotification(app, currConnection, message, webmodel.NOTE_JOIN_GROUP_REQUEST)
		if err != nil {
			return nil, err // handleJoinGroupNotification sends errors to wsconnection
		}

		err = app.DBModel.DeletUserFromGroup(note.GroupID.String, note.FromUserID.String)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("delete member user (id '%s') from group (id '%s') failed", note.FromUserID.String, note.GroupID.String), err)
		}

		newNote, err := createNotificationForGroup(app,
			note.FromUserID.String,
			currConnection.Client.UserID,
			webmodel.NOTE_JOIN_GROUP_REQUEST_DECLINED,
			"your request to join a group has been declined",
			note.GroupID.String,
		)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", note.FromUserID.String), err)
		}

		_, _, err = currConnection.SendMessageToUser(note.FromUserID.String, webmodel.NewNotification, newNote)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("sending notification to user (id'%s') failed", note.FromUserID.String), err)
		}

		return "declined", nil
	}
}

func ReplyAcceptInvitationToGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		note, err := handleJoinGroupNotification(app, currConnection, message, webmodel.NOTE_INVITE_TO_GROUP)
		if err != nil {
			return nil, err // handleJoinGroupNotification sends errors to wsconnection
		}

		err = app.DBModel.SetGroupMemberStatus(note.GroupID.String, note.UserID, true)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("set true member status for user (id '%s') in group (id '%s') failed", note.UserID, note.GroupID.String), err)
		}

		newNote, err := createNotificationForGroup(app,
			note.FromUserID.String,
			currConnection.Client.UserID,
			webmodel.NOTE_INVITE_TO_GROUP_ACCEPTED,
			"user has accepted the invitation to join a group",
			note.GroupID.String,
		)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", note.FromUserID.String), err)
		}

		_, _, err = currConnection.SendMessageToUser(note.FromUserID.String, webmodel.NewNotification, newNote)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("sending notification to user (id'%s') failed", note.FromUserID.String), err)
		}

		return "accepted", nil
	}
}

func ReplyDeclineInvitationToGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		note, err := handleJoinGroupNotification(app, currConnection, message, webmodel.NOTE_INVITE_TO_GROUP)
		if err != nil {
			return nil, err // handleJoinGroupNotification sends errors to wsconnection
		}

		err = app.DBModel.DeletUserFromGroup(note.GroupID.String, note.UserID)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("delete member user (id '%s') from group (id '%s') failed", note.UserID, note.GroupID.String), err)
		}

		newNote, err := createNotificationForGroup(app,
			note.FromUserID.String,
			currConnection.Client.UserID,
			webmodel.NOTE_INVITE_TO_GROUP_REJECTED,
			"user has rejected the invitation to join a group",
			note.GroupID.String,
		)
		if err != nil {
			return nil, handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", note.FromUserID.String), err)
		}

		_, _, err = currConnection.SendMessageToUser(note.FromUserID.String, webmodel.NewNotification, newNote)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("sending notification to user (id'%s') failed", note.FromUserID.String), err)
		}

		return "declined", nil
	}
}

func handleJoinGroupNotification(app *application.Application, currConnection *wsconnection.UsersConnection, message webmodel.WSMessage, noteType string) (*models.Notification, error) {
	notificationID, err := parse.PayloadToInt(message.Payload)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for notificationID: %s", message.Payload), err)
	}

	if notificationID < 1 {
		return nil, currConnection.WSBadRequest(message, "wrong notificationID")
	}
	note, err := app.DBModel.GetNotificationByID(notificationID)
	if err != nil {
		return nil, handleErrGetFromDBByID(app, currConnection, message, "notification", fmt.Sprintf("%d", notificationID), err)
	}

	if note.Type != noteType || !note.GroupID.Valid || !note.FromUserID.Valid || note.Read != 0 {
		return nil, currConnection.WSError(fmt.Sprintf("invalide notification for '%s', got '%s'", noteType, note), err)
	}

	err = app.DBModel.MarkNotificationRead(notificationID)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("mark notification (id '%d') as read failed", notificationID), err)
	}

	return note, nil
}
