package controllers

import (
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyCreateGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
	groupData, err := parse.PayloadToGroup(message.Payload)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a new group: %s", message.Payload), err)
	}

	errmessage := groupData.Validate()
	if errmessage != "" {
		return nil, currConnection.WSBadRequest(message, errmessage)
	}

	groupData.DateCreate = time.Now() // TODO get date from FrontEnd
	id, err := saveGroupToDB(app, currConnection, message, groupData)
	if err != nil {
		return nil, err
	}

	group := models.Group{
		ID:          id,
		Name:        groupData.Title,
		Description: groupData.Description,
		CreatorID:   currConnection.Client.UserID,
		DateCreate:  groupData.DateCreate,
		Members: []models.UserBase{
			*currConnection.Session.User,
		},
		Posts: nil,
	}

	return group, nil
} }

/*  sends the list of the group's members */
func ReplyGroupMembers(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
	groupID, err := parse.PayloadToString(message.Payload)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for groupID: %s", message.Payload), err)
	}

	if webmodel.IsEmpty(groupID) {
		return nil, currConnection.WSBadRequest(message, "groupID is empty")
	}

	members, err := app.DBModel.GetListOfGroupMembers(groupID)
	if err != nil {
		return nil, currConnection.WSError("getting list of members from DB failed", err)
	}

	return members, nil
} }

/* sends the list of groups which the current user belongs to */
func ReplyUserGroups(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
	userID, err := parse.PayloadToString(message.Payload)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for userID: %s", message.Payload), err)
	}
	if webmodel.IsEmpty(userID) {
		return nil, currConnection.WSBadRequest(message, "userID is empty")
	}
	groups, err := app.DBModel.GetListOfUserGroups(userID) // get list of groups
	if err != nil {
		return nil, currConnection.WSError("getting list of members from DB failed", err)
	}

	return groups, nil
} }

/* saves new group in DB and returns its ID */
func saveGroupToDB(app *application.Application, currConnection *wsconnection.UsersConnection, message webmodel.WSMessage, groupData webmodel.Group) (string, error) {
	var err error
	id, err := app.DBModel.CreateGroup(groupData.Title, groupData.Description, currConnection.Client.UserID, groupData.DateCreate)
	if err != nil {
		return "", handleErrAddToDB(app, currConnection, message,
			fmt.Sprintf("a new group '%s' created by userID '%s'", groupData.Title, currConnection.Client.UserID), err)
	}

	app.InfoLog.Printf("New Group '%s' is created. id: '%s'", groupData.Title, id)
	return id, nil
} 

func ReplyGetGroupProfile(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
	groupID, err := parse.PayloadToString(message.Payload)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a new group: %s", message.Payload), err)
	}

	if webmodel.IsEmpty(groupID) {
		return nil, currConnection.WSBadRequest(message, "groupID is empty")
	}

	group, err := app.DBModel.GetGroupByID(groupID)
	if err != nil {
		return nil, handleErrGetFromDBByID(app, currConnection, message, "group", groupID, err)
	}

	return group, nil
} }

func ReplyLeaveGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
	groupID, err := parse.PayloadToString(message.Payload)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a new group: %s", message.Payload), err)
	}

	if webmodel.IsEmpty(groupID) {
		return nil, currConnection.WSBadRequest(message, "groupID is empty")
	}

	err = app.DBModel.DeletUserFromGroup(groupID, currConnection.Client.UserID)
	if err != nil {
		return nil, currConnection.WSError(fmt.Sprintf("Could not get group from db with id: %s", message.Payload), err)
	}
	return "", nil
} }
