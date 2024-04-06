package controllers

import (
	"fmt"

	"social-network/backend/application"
	"social-network/backend/pkg/controllers/wsconnection"
	"social-network/backend/pkg/webmodel"
	"social-network/backend/pkg/webmodel/parse"
)

func ReplyAddCloseFriend(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		friendID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a new post: %s", message.Payload), err)
		}

		if webmodel.IsEmpty(friendID) {
			return nil, currConnection.WSBadRequest(message, "friendID is empty")
		}

		err = app.DBModel.AddCloseFriendToDB(friendID, currConnection.Client.UserID)
		if err != nil {
			return nil, currConnection.WSError("Couldn't add close friend to db", err)
		}

		return friendID, nil
	}
}
