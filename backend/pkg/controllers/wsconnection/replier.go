package wsconnection

import (
	"log"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wshub"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

type Replier interface {
	SendReply(currConnection *UsersConnection, wsMessage webmodel.WSMessage) error
}

type (
	FuncReplyCreator func(*UsersConnection, webmodel.WSMessage) (any, error)
	FuncReplier      func(*UsersConnection, webmodel.WSMessage) error
)

type WSmux struct {
	WShandlers      map[string]Replier
	InfoLog, ErrLog *log.Logger
	Hub             *wshub.Hub
}

func (fRC FuncReplyCreator) SendReply(currConnection *UsersConnection, wsMessage webmodel.WSMessage) error {
	replyData, err := fRC(currConnection, wsMessage)
	if err != nil {
		return err
	}

	return currConnection.SendReply(wsMessage, replyData)
}

func (fR FuncReplier) SendReply(currConnection *UsersConnection, wsMessage webmodel.WSMessage) error {
	return fR(currConnection, wsMessage)
}
