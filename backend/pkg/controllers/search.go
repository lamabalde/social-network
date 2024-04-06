package controllers

import (
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplySearchGroupsUsers(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		searchQuery, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid query '%s'", message.Payload), err)
		}
		results, err := app.DBModel.GetListOfGroupsUsers(searchQuery)
		if err != nil {
			return nil, currConnection.WSError("Failed to get search data from database", err)
		}
		return results, nil
	}
}

func ReplySearchUsersNotFriends(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		searchQuery, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for GroupUser  '%s'", message.Payload), err)
		}
		userID := currConnection.Client.UserID

		results, err := app.DBModel.GetListOfUsersNotFriends(searchQuery, userID)
		if err != nil {
			return nil, currConnection.WSError("Failed to get search data from database", err)
		}

		return results, nil
	}
}

func ReplySearchUsersNotInGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		searchQuery, err := parse.PayloadToGroupUser(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("ReplySearchUsersNotInGroup: Invalid payload for GroupUser  '%s'", message.Payload), err)
		}

		errmessage := searchQuery.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		results, err := app.DBModel.GetUsersByPartialNameNotInGroup(searchQuery.User, searchQuery.GroupID)
		if err != nil {
			return nil, currConnection.WSError("ReplySearchUsersNotInGroup: Failed to get search data from database", err)
		}

		return results, nil
	}
}
