package controllers

import (
	"errors"
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyCreateGroupEvent(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		eventData, err := parse.PayloadToGroupEvent(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a new group: %s", message.Payload), err)
		}

		errmessage := eventData.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		event := createEventFromWebModelEvent(currConnection, eventData)

		id, err := app.DBModel.CreateEvent(event)
		if err == models.ErrNoRecords {
			return nil, currConnection.WSBadRequest(message, fmt.Sprintf("can not create event: no group '%s' or  user '%s' ", event.GroupID, event.CreatorID))
		}
		if err != nil {
			return "", handleErrAddToDB(app, currConnection, message,
				fmt.Sprintf("a new group '%s' created by userID '%s'", eventData.Title, currConnection.Client.UserID), err)
		}

		app.InfoLog.Printf("New Event '%s' is created. id: '%s'", eventData.Title, id)

		event.ID = id

		return event, nil
	}
}

func createEventFromWebModelEvent(currConnection *wsconnection.UsersConnection, eventData webmodel.Event) models.Event {
	return models.Event{
		Title:       eventData.Title,
		Description: eventData.Description,
		GroupID:     eventData.GroupID,
		CreatorID:   currConnection.Client.UserID,
		DateEvent:   eventData.DateEvent,
		DateCreate:  time.Now(),
	}
}

func ReplySetUserOptionForEvent(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		eventOption, err := parse.PayloadToUserOptionForEvent(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for UserOptionForEvent: %s", message.Payload), err)
		}

		errmessage := eventOption.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		id, err := app.DBModel.AddEventMember(eventOption.EventID, currConnection.Client.UserID, eventOption.Option)
		if errors.Is(err, models.ErrNoRecords) {
			return nil, currConnection.WSBadRequest(message, fmt.Sprintf("can not add user option: wrong event id'%s' or  user '%s' is not a member of event's group", eventOption.EventID, currConnection.Client.UserID))
		}
		if err != nil {
			return "", handleErrAddToDB(app, currConnection, message,
				fmt.Sprintf("a new member '%s' of event '%s'", currConnection.Client.UserID, eventOption.EventID), err)
		}

		app.InfoLog.Printf("New member '%s' of event '%s' is added. id: '%d'", currConnection.Client.UserID, eventOption.EventID, id)

		return eventOption, nil
	}
}

func ReplyChangeUserOptionForEvent(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		eventOption, err := parse.PayloadToUserOptionForEvent(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for UserOptionForEvent: %s", message.Payload), err)
		}

		errmessage := eventOption.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		err = app.DBModel.ChangeEventMemberOption(eventOption.EventID, currConnection.Client.UserID, eventOption.Option)
		if errors.Is(err, models.ErrNoRecords) {
			return nil, currConnection.WSBadRequest(message, fmt.Sprintf("can not change user option: wrong event id'%s' or  user '%s' is not a member of event's group", eventOption.EventID, currConnection.Client.UserID))
		}
		if err != nil {
			return "", currConnection.WSError(fmt.Sprintf("changing user's (id'%s') option for event '%s' failed", currConnection.Client.UserID, eventOption.EventID), err)
		}

		app.InfoLog.Printf("user's (id '%s') option for event '%s' is changed.", currConnection.Client.UserID, eventOption.EventID)

		return eventOption, nil
	}
}

func ReplyGetGroupEvent(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		eventID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a eventID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(eventID) {
			return nil, currConnection.WSBadRequest(message, "eventID is empty")
		}

		event, err := app.DBModel.GetEventByID(eventID)
		if err != nil {
			return nil, handleErrGetFromDBByID(app, currConnection, message, "event", eventID, err)
		}

		return event, nil
	}
}

func ReplyGetGroupEventsList(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		groupID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for groupID: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(groupID) {
			return nil, currConnection.WSBadRequest(message, "userID is empty")
		}

		events, err := app.DBModel.GetListOfGroupEvents(groupID)
		if err != nil {
			return nil, currConnection.WSError("getting list of events from DB failed", err)
		}

		return events, nil
	}
}
