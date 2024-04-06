package controllers

import (
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyFullPostAndComments(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		postID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid postID '%s'", message.Payload), err)
		}

		if webmodel.IsEmpty(postID) {
			return nil, currConnection.WSBadRequest(message, "postID is empty")
		}

		post, err := getPost(app, currConnection, postID, message)
		if err != nil {
			return nil, err
		}
		return post, nil
	}
}
