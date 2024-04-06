package controllers

import (
	"fmt"

	"social-network/backend/application"
	"social-network/backend/pkg/controllers/wsconnection"
	"social-network/backend/pkg/webmodel"
	"social-network/backend/pkg/webmodel/parse"
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
