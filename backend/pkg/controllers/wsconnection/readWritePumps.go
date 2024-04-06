package wsconnection

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

var NewLine = []byte("\n")

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (uc *UsersConnection) ReadPump(app *application.Application) {
	defer func() {
		// (online changes)uc.WsServer.Hub.UnRegisterFromHub(uc.Client)
		err := uc.deleteClientAndSendUserOffline(app)
		if err != nil {
			uc.WsServer.ErrLog.Printf("ReadPump: error client delete: %v", err)
		}

		close(uc.Client.ReceivedMessages)
		err = uc.Client.Conn.Close()
		if err != nil {
			uc.WsServer.ErrLog.Printf("ReadPump: error closing connection %p: %v", uc.Client.Conn, err)
		}
		uc.WsServer.InfoLog.Printf("ReadPump closed connection %p", uc.Client.Conn)
	}()

	uc.Client.Conn.SetReadLimit(maxMessageSize)
	uc.Client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	uc.Client.Conn.SetPongHandler(func(string) error { uc.Client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message webmodel.WSMessage
		err := uc.Client.Conn.ReadJSON(&message)
		fmt.Printf("Message received from js: %s\n ", message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				uc.WsServer.ErrLog.Printf("websocket connection %p to '%s' was unexpected closed: %#v", uc.Client.Conn, uc.Client.Conn.LocalAddr(), err)
			}
			uc.WsServer.ErrLog.Printf("ReadPump is closing connection %p  of client  '%s' : %#v", uc.Client.Conn, uc.Client, err)
			break
		}

		replier, ok := uc.WsServer.WShandlers[message.Type]
		if !ok {
			uc.WsServer.ErrLog.Printf("unknown type message received: %s", message.Type)
			continue
		}

		err = replier.SendReply(uc, message)
		if err != nil && !errors.Is(err, webmodel.ErrWarning) {
			break
		}

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (uc *UsersConnection) WritePump(app *application.Application) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := uc.Client.Conn.Close()
		if err != nil {
			uc.WsServer.ErrLog.Printf("WritePump: error closing connection: %v", err)
		} else {
			uc.WsServer.InfoLog.Printf("WritePump closed connection %p", uc.Client.Conn)
		}
	}()
	for {
		chann := uc.Client.ReceivedMessages
		select {
		case message, ok := <-chann:
			// fmt.Printf("write message %s\n", message)
			if !ok {
				// The hub closed the channel.
				uc.Client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				uc.WsServer.InfoLog.Printf("WritePump is closing connection because the hub closed the channel %p ", chann)
				uc.WsServer.InfoLog.Printf("WritePump is closing connection %p of client '%s' because the hub closed the channel %p ", uc.Client.Conn, uc.Client, uc.Client.ReceivedMessages)
				return
			}

			uc.Client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := uc.Client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				uc.WsServer.ErrLog.Printf("cannot create the NextWriter on the connenction %p : %v", uc.Client.Conn, err)
				return
			}
			uc.writeMessage(w, message)

			// Add queued chat messages to the current websocket message.
			n := len(chann)
			for i := 0; i < n; i++ {
				message = <-chann
				uc.writeMessage(w, message)
			}

			if err := w.Close(); err != nil {
				uc.WsServer.ErrLog.Printf("cannot close the writer on the connenction %p : %v", uc.Client.Conn, err)
				return
			}
		case <-ticker.C:
			uc.Client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := uc.Client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				uc.WsServer.ErrLog.Printf("ping the connenction %s failed: %v", uc.Client.Conn.LocalAddr(), err)
				return
			}
		}
	}
}

func (uc *UsersConnection) writeMessage(w io.WriteCloser, message []byte) error {
	_, err := w.Write(message)
	if err != nil {
		return err
	}
	uc.WsServer.InfoLog.Printf("Websocket: send message: '%s' to client %s", helpers.ShortMessage(message), uc.Client)
	_, err = w.Write(NewLine)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UsersConnection) deleteClientAndSendUserOffline(app *application.Application) error {
	uc.WsServer.Hub.UnRegisterClientFromHub(uc.Client)
	if uc.Client.UserID != "" {
		roomID, err := json.Marshal(uc.Client.Room.ID)
		if err != nil {
			return err
		}
		return uc.WsServer.WShandlers[webmodel.UserOffline].SendReply(uc, webmodel.WSMessage{Type: webmodel.UserOffline, Payload: roomID})
	}
	return nil
}
