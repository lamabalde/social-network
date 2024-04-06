package controllers

import (
	"database/sql"
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func createNotificationForUser(app *application.Application, forUserID, noteType, body string) (*models.Notification, error) {
	note := &models.Notification{
		UserID:     forUserID,
		Type:       noteType,
		Body:       body,
		Read:       0,
		DateCreate: time.Now(),
	}

	id, err := app.DBModel.AddNotification(note)
	if err != nil {
		return nil, fmt.Errorf("createNotificationFromUser: creating notification for user (id'%s') failed: %w", forUserID, err)
	}

	note.ID = id

	app.InfoLog.Printf("create a notification '%s' for user (id'%s')", noteType, forUserID)
	return note, nil
}

func createNotificationFromUser(app *application.Application, forUserID, fromUserID, noteType, body string) (*models.Notification, error) {
	note := &models.Notification{
		UserID:     forUserID,
		Type:       noteType,
		Body:       body,
		FromUserID: sql.NullString{String: fromUserID, Valid: true},
		Read:       0,
		DateCreate: time.Now(),
	}

	id, err := app.DBModel.AddNotification(note)
	if err != nil {
		return nil, fmt.Errorf("createNotificationFromUser: creating notification for user (id'%s') failed: %w", forUserID, err)
	}

	note.ID = id

	app.InfoLog.Printf("create a notification '%s' from user (id'%s') for user (id'%s')", noteType, fromUserID, forUserID)
	return note, nil
}

func createNotificationForGroup(app *application.Application, forUserID, fromUserID, noteType, body, groupID string) (*models.Notification, error) {
	note := &models.Notification{
		UserID:     forUserID,
		Type:       noteType,
		Body:       body,
		FromUserID: sql.NullString{String: fromUserID, Valid: true},
		GroupID:    sql.NullString{String: groupID, Valid: true},
		Read:       0,
		DateCreate: time.Now(),
	}

	id, err := app.DBModel.AddNotification(note)
	if err != nil {
		return nil, fmt.Errorf("createNotificationForGroup: creating notification for user (id'%s') failed: %w", forUserID, err)
	}

	note.ID = id

	app.InfoLog.Printf("create a notification '%s' for user (id'%s')", noteType, forUserID)
	return note, nil
}

func ReplyGetUserNotifications(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(userID) {
			return nil, currConnection.WSBadRequest(message, "groupID is empty")
		}

		notifications, err := app.DBModel.GetUserUnReadNotification(userID)
		if err != nil {
			return nil, fmt.Errorf("getting notifications for user (id'%s') failed", userID)
		}

		return notifications, nil
	}
}
