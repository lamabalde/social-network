package wsconnection

import (
	"encoding/json"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wshub"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/session"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

type UsersConnection struct {
	Session  *session.Session
	Client   *wshub.Client
	WsServer WSmux
}



/*
sends uc.Client a successful reply to the requestMessage with the data as payload
*/
func (uc *UsersConnection) SendReply(requestMessage webmodel.WSMessage, data any) error {
	wsMessage, err := requestMessage.CreateReplyToRequestMessage("success", data)
	if err != nil {
		return uc.WSErrCreateMessage(err)
	}

	uc.Client.WriteMessage(wsMessage)
	uc.WsServer.InfoLog.Printf("send message %s to channel of client %p", helpers.ShortMessage(wsMessage), uc.Client)
	return nil
}

/*
sends uc.Client a successful message with the type = 'messageType' and 'data' as payload.
It returns the message converted into json.
It returns error if the message could not be converted to json.
*/
func (uc *UsersConnection) SendSuccessMessage(messageType string, data any) (json.RawMessage, error) {
	wsMessage, err := webmodel.CreateJSONMessage(messageType, "success", data)
	if err != nil {
		return wsMessage, uc.WSErrCreateMessage(err)
	}

	uc.Client.WriteMessage(wsMessage)
	return wsMessage, nil
}

/*
sends 'recipient' client a successful message with the type = 'messageType' and 'data' as payload.
It returns the message converted into json.
It returns error if the message could not be converted to json.
*/
func (uc *UsersConnection) SendMessageToOtherClient(recipient *wshub.Client, messageType string, data any) (json.RawMessage, error) {
	wsMessage, err := webmodel.CreateJSONMessage(messageType, "success", data)
	if err != nil {
		return wsMessage, uc.WSErrCreateMessage(err)
	}

	recipient.WriteMessage(wsMessage)
	return wsMessage, nil
}

/*
sends a client of the user with given 'userID' a successful message with the type = 'messageType' and 'data' as payload.
It returns the message converted into json and true if the message was sent successfully.
It returns error if the message could not be converted to json.
*/
func (uc *UsersConnection) SendMessageToUser(userID string, messageType string, data any) (json.RawMessage, bool, error) {
	wsMessage, err := webmodel.CreateJSONMessage(messageType, "success", data)
	if err != nil {
		return wsMessage, false, uc.WSErrCreateMessage(err)
	}

	ok := uc.SendBytesToUser(userID, wsMessage)

	return wsMessage, ok, nil
}

/*
sends bytes to a client of the user with given 'userID'.
It returns the message converted into json and true if the message was sent successfully.
It returns error if the message could not be converted to json.
*/
func (uc *UsersConnection) SendBytesToUser(userID string, rawData []byte) bool {
	sentMark := uc.WsServer.Hub.SendMessageToUsers(rawData, []string{userID})
	_, ok := sentMark[userID]

	return ok
}

/*
sends clients of the users on the list a successful message with the type = 'messageType' and 'data' as payload.
It returns the message converted into json and true if the message was sent successfully.
It returns error if the message could not be converted to json.
*/
func (uc *UsersConnection) SendMessageToUsers(usersID []string, messageType string, data any) (json.RawMessage, map[string]bool, error) {
	wsMessage, err := webmodel.CreateJSONMessage(messageType, "success", data)
	if err != nil {
		return wsMessage, nil, uc.WSErrCreateMessage(err)
	}

	sentMark := uc.SendBytesToUsers(usersID, wsMessage)

	return wsMessage, sentMark, nil
}

/*
sends bytes to clients of the users on the list.
It returns the message converted into json and true if the message was sent successfully.
It returns error if the message could not be converted to json.
*/
func (uc *UsersConnection) SendBytesToUsers(usersID []string, rawData []byte) map[string]bool {
	sentMark := uc.WsServer.Hub.SendMessageToUsers(rawData, usersID)

	return sentMark
}

/*
sends to the 'room' a successful message with the type = 'messageType' and 'data' as payload.
It returns the message converted into json and a map with true value if the message was sent successfully and false otherwise.
It returns error if the message could not be converted to json.
*/
func (uc *UsersConnection) SendMessageToClientRoom(messageType string, data any) (json.RawMessage, map[string]bool, error) {
	wsMessage, err := webmodel.CreateJSONMessage(messageType, "success", data)
	if err != nil {
		return wsMessage, nil, uc.WSErrCreateMessage(err)
	}

	sentMarksMap := uc.SendBytesToClientRoom(wsMessage)

	return wsMessage, sentMarksMap, nil
}

/*
sends bytes to the 'room' .
It returns the message converted into json and a map with true value if the message was sent successfully and false otherwise.
*/
func (uc *UsersConnection) SendBytesToClientRoom(rawData []byte) map[string]bool {
	sentMarksMap := uc.WsServer.Hub.BroadcastMessageInRoom(rawData, uc.Client.Room)

	return sentMarksMap
}
