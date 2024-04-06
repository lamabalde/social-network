package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wshub"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/errorhandle"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/session"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

const (
	CREATE_CHAT_URL = "/openChat/"
	JOIN_CHAT_URL   = "/joinChat/"
)

// IndexWs handles websocket requests from the main page. URL: /ws
func IndexWs(app *application.Application, wsReplyersSet wsconnection.WSmux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		fmt.Println("IndexWs url: ", r.URL.Path)
		sess, ok := r.Context().Value(CTX_USER).(*session.Session)
		if !ok {
			errorhandle.Forbidden(app, w, r, "error getting session from request context")
			return
		}

		if !sess.IsLoggedin() {
			errorhandle.Forbidden(app, w, r, "unauthorized session in request context")
			return
		}
		fmt.Printf("IndexWs runs with session: %s\n", sess)

		conn, err := app.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			errorhandle.ServerError(app, w, r, "Upgrade failed:", err)
			return
		}
		app.InfoLog.Printf("connection %p to '%s' is upgraded to the WebSocket protocol", conn, r.URL.Path)
		// the connection will be closed in WritePump or ReadPump functions

		client, ok := wshub.NewClient(app.Hub, sess.User.ID, sess.User.UserName, app.Hub.UniRoom, conn, nil, nil)
		if !ok {
			errorhandle.ServerError(app, w, r, "cant create a client in Universal room", errors.New("universal room does not exist"))
			return
		}
		currentConnection := &wsconnection.UsersConnection{
			Session:  sess,
			Client:   client,
			WsServer: wsReplyersSet,
		}

		go currentConnection.WritePump(app)
		go currentConnection.ReadPump(app)

		if currentConnection.Session.IsLoggedin() {
			err = controllers.SendOnlineUsers(app, currentConnection)
			if err != nil && !errors.Is(err, webmodel.ErrWarning) {
				logErrorAndCloseConn(app, conn, "send online users list failed", err)
				return
			}
		}
		app.InfoLog.Printf("registered new client: %s", currentConnection.Client)

		err = sendSession(app, currentConnection)
		if err != nil && !errors.Is(err, webmodel.ErrWarning) {
			logErrorAndCloseConn(app, conn, "send session failed", err)
			return
		}
	}
}

func sendSession(app *application.Application, currConn *wsconnection.UsersConnection) error {
	_, err := currConn.SendSuccessMessage(webmodel.CurrentSession, currConn.Session)
	return err
}

func logErrorAndCloseConn(app *application.Application, conn *websocket.Conn, errMessage string, err error) {
	app.ErrLog.Printf("%s: %v", errMessage, err)

	err = conn.Close()
	if err != nil {
		app.ErrLog.Printf("error closing connection: %v", err)
	}
}

func OpenChat(app *application.Application, wsReplyersSet wsconnection.WSmux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create chat and join current user to it
		var err error
		sess, ok := r.Context().Value(CTX_USER).(*session.Session)
		if !ok {
			errorhandle.Forbidden(app, w, r, "error getting session from request context")
			return
		}
		if !sess.IsLoggedin() {
			errorhandle.Forbidden(app, w, r, "unauthorized session in request context")
			return
		}

		chatType, id, err := getChatParams(r.URL)
		if err != nil {
			errorhandle.BadRequestError(app, w, r, err.Error())
			return
		}

		chat, ok := getChat(app, w, r, sess.User.ID, id, chatType)
		if !ok {
			return
		}
		chatRoom, err := getChatRoom(app.Hub, chat.ID, chatType)
		if err != nil {
			errorhandle.BadRequestError(app, w, r, err.Error())
			return
		}

		conn, err := app.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			errorhandle.ServerError(app, w, r, "Upgrade failed:", err)
			return
		}
		app.InfoLog.Printf("connection %p to '%s' is upgraded to the WebSocket protocol", conn, r.URL.Path)

		client, ok := wshub.NewClient(app.Hub, sess.User.ID, sess.User.UserName, chatRoom, conn, nil, nil)
		if !ok {
			errorhandle.ServerError(app, w, r, "cant create a client:", fmt.Errorf("room id '%s' does not exist", chatRoom.ID))
			return
		}

		currentConnection := &wsconnection.UsersConnection{
			Session:  sess,
			Client:   client,
			WsServer: wsReplyersSet,
		}
		// the connection will be closed in WritePump or ReadPump functions
		app.InfoLog.Printf("New client in room '%s' is created: %s", chatRoom, currentConnection.Client)

		go currentConnection.WritePump(app)

		go currentConnection.ReadPump(app)

		err = controllers.SendChattingUsers(app, currentConnection)
		if err != nil && !errors.Is(err, webmodel.ErrWarning) {
			logErrorAndCloseConn(app, conn, "sending chatting users failed", err)
			return
		}

		app.InfoLog.Printf("User '%s' opened %s chat id '%s'.", sess.User, chatType, chat.ID)
	}
}

/*
gets type of chat and the id from the request's url.
URL must be  "/openChat/private/?id=<recepientID>" or "/openChat/group/?id=<groupID>"
*/
func getChatParams(url *url.URL) (chatType, id string, err error) {
	var ok bool
	chatType, ok = strings.CutPrefix(url.Path, CREATE_CHAT_URL)
	if !ok {
		err = fmt.Errorf("wrong url '%s' in CreateChat Handler", url.Path)
		return
	}

	chatType = strings.ToLower(chatType)

	id = url.Query().Get("id")
	if id == "" {
		err = fmt.Errorf("cant open the chat: no id specified")
		return
	}

	return
}

/*
gets from DB(creates if necessary) chat and returns the chat object and  true if it succeeded, false otherwise.
If it failed sends an error response to w.
*/
func getChat(app *application.Application, w http.ResponseWriter, r *http.Request, currentUserID, goalID, chatType string) (models.Chat, bool) {
	var chat models.Chat
	var err error
	switch chatType {
	case controllers.PRIVATE_CHAT_ROOM:
		chat, err = controllers.OpenPrivateChat(app, currentUserID, goalID)
		if err != nil {
			errorhandle.ServerError(app, w, r, "error creating private chat", err)
			return chat, false
		}
	case controllers.GROUP_CHAT_ROOM:
		chat, err = controllers.OpenGroupChat(app, goalID)
		if err != nil {
			errorhandle.ServerError(app, w, r, "error creating group chat", err)
			return chat, false
		}
	default:
		errorhandle.BadRequestError(app, w, r, "wrong url for opening chat")
		return chat, false
	}
	return chat, true
}

/*
for a private chat gets a room from hub,create a new one if it doesn't exist.
for a group chat creates a new room , returns an error, if it already exists.
*/
func getChatRoom(hub *wshub.Hub, chatID string, chatType string) (*wshub.Room, error) {
	chatRoom, ok := hub.GetRoom(chatID)
	if !ok {
		chatRoom, ok = wshub.NewRoom(hub, chatID, chatType)
		if !ok {
			return nil, fmt.Errorf("private chat id '%s' was already created, try again", chatID)
		}
	}
	// TODO delete it (trying to rid of joinChat (that block was used with joinChat))
	// else {
	// 	if chatType == GROUP_CHAT {
	// 		return nil, fmt.Errorf("group chat for group id '%s' is already open", chatID)
	// 	}
	// }
	return chatRoom, nil
}

// func JoinGroupChat(app *application.Application, wsReplyersSet wsconnection.WSmux) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// create chat and join current user to it
// 		var err error

// 		sess, ok := r.Context().Value(CTX_USER).(*session.Session)
// 		if !ok {
// 			errorhandle.Forbidden(app, w, r, "error getting session from request context")
// 			return
// 		}
// 		if !sess.IsLoggedin() {
// 			errorhandle.Forbidden(app, w, r, "unauthorized session in request context")
// 			return
// 		}

// 		groupID := r.URL.Query().Get("id")
// 		if groupID == "" {
// 			errorhandle.BadRequestError(app, w, r, "cant join to chat: no id specified for group")
// 			return
// 		}

// 		chatRoom, ok := app.Hub.GetRoom(groupID)
// 		if !ok {
// 			errorhandle.BadRequestError(app, w, r, fmt.Sprintf("cant join to chat: room id '%s' is not registred in the Hub", groupID))
// 			return
// 		}

// 		conn, err := app.Upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			errorhandle.ServerError(app, w, r, "Upgrade failed:", err)
// 			return
// 		}
// 		app.InfoLog.Printf("connection %p to '%s' is upgraded to the WebSocket protocol", conn, r.URL.Path)

// 		client, ok := wshub.NewClient(app.Hub, sess.User.ID, sess.User.UserName, chatRoom, conn, nil, nil)
// 		if !ok {
// 			errorhandle.ServerError(app, w, r, "cant create a client:", fmt.Errorf("room id '%s' does not exist", groupID))
// 			return
// 		}
// 		currentConnection := &wsconnection.UsersConnection{
// 			Session:  sess,
// 			Client:   client,
// 			Repliers: wsReplyersSet,
// 		}
// 		// the connection will be closed in WritePump or ReadPump functions
// 		app.InfoLog.Printf("New client in room '%s' is created: %s", chatRoom, currentConnection.Client)

// 		go currentConnection.WritePump(app)

// 		go currentConnection.ReadPump(app)

// 		app.InfoLog.Printf("User '%s' joined to chat id '%s'.", sess.User, chatRoom.ID)

// 		// TODO send message to the room {"userJoinChat" userID}
// 	}
// }
