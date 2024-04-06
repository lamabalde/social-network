package controllers

import (
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyUserFollowers(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(userID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		followers, err := app.DBModel.GetFollowers(userID)
		if err != nil {
			return nil, currConnection.WSError("getting list of followers from DB failed", err)
		}

		return followers, nil
	}
}

func ReplyUserFollowing(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(userID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		followings, err := app.DBModel.GetFollowing(userID)
		if err != nil {
			return nil, currConnection.WSError("getting list of followings from DB failed", err)
		}

		return followings, nil
	}
}

func ReplyFollowUser(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		var followStatus string
		userID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(userID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		user, err := app.DBModel.GetUserByID(userID)
		if err != nil {
			return nil, handleErrGetFromDBByID(app, currConnection, message, "user", userID, err)
		}

		if user.ProfileType == webmodel.PRIVATE {
			followStatus = models.FOLLOW_STATUS_REQUESTED
			err = app.DBModel.AddFollowing(currConnection.Client.UserID, userID, followStatus)
			if err != nil {
				return nil, handleErrAddToDB(app, currConnection, message, "following", err)
			}

			note, err := createNotificationFromUser(app, userID, currConnection.Client.UserID,
				webmodel.NOTE_FOLLOW_REQUEST, fmt.Sprintf("%s wants to follow you", currConnection.Client.UserName))
			if err != nil {
				return nil, handleErrAddToDB(app, currConnection, message, fmt.Sprintf("notification for user (id'%s')", userID), err)
			}

			_, _, err = currConnection.SendMessageToUser(userID, webmodel.NewNotification, note)
			if err != nil {
				return nil, currConnection.WSError(fmt.Sprintf("sending notification to user (id'%s') failed", userID), err)
			}

			return followStatus, nil
		}

		followStatus = models.FOLLOW_STATUS_FOLLOWING
		err = app.DBModel.AddFollowing(currConnection.Client.UserID, userID, followStatus)
		if err != nil {
			return nil, currConnection.WSError("insert following to DB failed", err)
		}

		reply := webmodel.FollowingReply{Id: userID, FollowStatus: followStatus}

		return reply, nil
	}
}

func ReplyGetFollowStatus(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(userID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		followStatus, err := app.DBModel.GetFollowStatus(currConnection.Client.UserID, userID)
		if err != nil {
			return nil, currConnection.WSError("getting list of followings from DB failed", err)
		}
		return followStatus, nil
	}
}

func ReplyAcceptFollowRequest(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		// userID, err := parse.PayloadToString(message.Payload)
		followResponse, err := parse.PayloadToFollowResponseWithNotificationID(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(followResponse.UserID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		err = app.DBModel.SetFollowStatus(followResponse.UserID, currConnection.Client.UserID, models.FOLLOW_STATUS_FOLLOWING)
		if err == models.ErrNoRecords {
			return nil, currConnection.WSBadRequest(message, fmt.Sprintf("no request to follow from user '%s' to user '%s'", followResponse.UserID, currConnection.Client.UserID))
		}
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("set follow status from user '%s' to user '%s' in DB failed", followResponse.UserID, currConnection.Client.UserID), err)
		}

		err = app.DBModel.MarkNotificationRead(followResponse.Id)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("mark notification (id '%d') as read failed", followResponse.Id), err)
		}

		return followResponse, nil
	}
}

func ReplyDeclineFollowRequest(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		followResponse, err := parse.PayloadToFollowResponseWithNotificationID(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}
		if webmodel.IsEmpty(followResponse.UserID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		err = app.DBModel.DeleteFollowing(followResponse.UserID, currConnection.Client.UserID)
		if err == models.ErrNoRecords {
			return nil, currConnection.WSBadRequest(message, fmt.Sprintf("no request to follow from user '%s' to user '%s'", followResponse.UserID, currConnection.Client.UserID))
		}
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("set follow status from user '%s' to user '%s' in DB failed", followResponse.UserID, currConnection.Client.UserID), err)
		}

		err = app.DBModel.MarkNotificationRead(followResponse.Id)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("mark notification (id '%d') as read failed", followResponse.Id), err)
		}

		return followResponse, nil
	}
}

func ReplyUnFollowUser(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(userID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		err = app.DBModel.DeleteFollowing(currConnection.Client.UserID, userID)
		if err == models.ErrNoRecords {
			return nil, currConnection.WSBadRequest(message, fmt.Sprintf("no request to follow from user '%s' to user '%s'", currConnection.Client.UserID, userID))
		}
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("set follow status from user '%s' to user '%s' in DB failed", currConnection.Client.UserID, userID), err)
		}

		return userID, nil
	}
}
