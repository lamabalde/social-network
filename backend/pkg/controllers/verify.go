package controllers

import (
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyVerifyGroupView(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		groupID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid groupID '%s'", message.Payload), err)
		}

		if webmodel.IsEmpty(groupID) {
			return nil, currConnection.WSBadRequest(message, "groupID is empty")
		}
		userID := currConnection.Client.UserID

		memberStatus, err := app.DBModel.CheckGroupMemberStatus(groupID, userID)
		if err != nil {
			app.ErrLog.Printf("checking group member status failed: %v", err)
			return false, err
		}
		// check if status is requested or joined
		return memberStatus, nil
	}
}
