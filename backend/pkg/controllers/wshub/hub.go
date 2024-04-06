package wshub

import (
	"encoding/json"
	"fmt"
)

const (
	UNIVERSAL_ROOM_ID = "UniversalRoom"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered rooms.
	rooms   *SafeRoomsMap
	UniRoom *Room

	// Inbound messages from a client.
	broadcast chan *message

	sendToListOfUsers chan *messageToListOfUsers

	clientRegister   chan *Client
	clientUnregister chan *Client

	roomRegister   chan *Room
	roomUnregister chan *Room

	// OnlineUsersRequest chan *Client
}

func NewHub() *Hub {
	hub := &Hub{
		broadcast:         make(chan *message),
		sendToListOfUsers: make(chan *messageToListOfUsers),
		rooms:             NewSafeRoomsMap(),
		clientRegister:    make(chan *Client),
		clientUnregister:  make(chan *Client),
		roomRegister:      make(chan *Room),
		roomUnregister:    make(chan *Room),
		// OnlineUsersRequest: make(chan *Client),
	}
	uniRoom := createRoom(hub, UNIVERSAL_ROOM_ID, UNIVERSAL_ROOM_ID)
	hub.rooms.Set(UNIVERSAL_ROOM_ID, uniRoom)
	hub.UniRoom = uniRoom

	return hub
}

type message struct {
	content   json.RawMessage
	room      *Room
	usersSentTo chan map[string]bool
}

type messageToListOfUsers struct {
	content json.RawMessage
	usersID []string
	usersSentTo chan map[string]bool
}

func (h *Hub) Run() {
	for { // TODO does it need to use locking here? (the other hub methods use locking)
		select {
		case room := <-h.roomRegister:
			if !h.isThereRoom(room) {
				h.rooms.Set(room.ID, room)
				room.Registered <- true
			} else {
				fmt.Printf("room id %s is already registered", room.ID)
				room.Registered <- false
			}

		case room := <-h.roomUnregister:
			if h.isThereRoom(room) {
				h.rooms.Delete(room.ID)
			}

		case client := <-h.clientRegister:
			if h.isThereRoom(client.Room) {
				client.Room.Clients.Set(client.UserID, client)

				client.Registered <- true
			} else {
				client.Registered <- false
			}

		case client := <-h.clientUnregister:
			if h.isThereRoom(client.Room) {
				if client.Room.isThereClient(client) {
					client.Room.DeleteClient(client)
					// TODO: if client.Room.Clients.Len() != 0 send a message fmt.Sprintf("user %s left the chat", client.UserName)
				}

				if client.Room.ID != UNIVERSAL_ROOM_ID && client.Room.Clients.Len() == 0 {
					h.rooms.Delete(client.Room.ID)
				}
			}

		case message := <-h.broadcast:
			// h.rooms.Lock()

			if h.isThereRoom(message.room) {
				usersSentTo := make(map[string]bool)
				for _, client := range message.room.Clients.items {
					// select {
					// case client.ReceivedMessages <- message.content:
					// default:
					// 	// If the client's send buffer is full, then the hub assumes that the client is dead or stuck.
					// 	// In this case, the hub unregisters the client.
					// 	close(client.ReceivedMessages)
					// 	message.room.DeleteClient(client)
					// }
					sentMessageToClient(message.content, client)

					usersSentTo[client.UserID] = true
				}

				message.usersSentTo <- usersSentTo

				if message.room.Clients.Len() == 0 {
					h.rooms.Delete(message.room.ID)
				}
			}
		case message := <-h.sendToListOfUsers:
			// h.rooms.Lock()

			usersSentTo := make(map[string]bool)

			for _,userID := range message.usersID {
				if client, ok := h.UniRoom.Clients.Get(userID); ok {
					// select {
					// case client.ReceivedMessages <- message.content:
					// default:
					// 	// If the client's send buffer is full, then the hub assumes that the client is dead or stuck.
					// 	// In this case, the hub unregisters the client.
					// 	close(client.ReceivedMessages)
					// 	h.UniRoom.DeleteClient(client)
					// }
					sentMessageToClient(message.content, client)

					usersSentTo[client.UserID] = true
				}
			}

			message.usersSentTo <- usersSentTo

			// h.rooms.Unlock()
			// case client := <-h.OnlineUsersRequest:

		}
	}
}

func sentMessageToClient(messageContent json.RawMessage, client *Client) {
	select {
	case client.ReceivedMessages <- messageContent:
	default:
		// If the client's send buffer is full, then the hub assumes that the client is dead or stuck.
		// In this case, the hub unregisters the client.
		close(client.ReceivedMessages)
		client.Room.DeleteClient(client)
	}
}

type MapID map[string]*Client

func (m MapID) CheckID(userID string) bool {
	_, ok := m[userID]
	return ok
}

// RegisterRoomToHub registers the room to its hub
func (h *Hub) RegisterRoomToHub(r *Room) {
	h.roomRegister <- r
}

// UnRegisterRoomToHub removes the room from its hub
func (h *Hub) UnRegisterRoomFromHub(r *Room) {
	h.roomUnregister <- r
}

// RegisterClientToHub registers the client to its hub
func (h *Hub) RegisterClientToHub(c *Client) {
	h.clientRegister <- c
}

// UnRegisterClientToHub removes the client from its hub
func (h *Hub) UnRegisterClientFromHub(c *Client) {
	h.clientUnregister <- c
}

func (h *Hub) GetOnlineUsers() MapID {
	usersID := make(MapID, h.UniRoom.Clients.Len())
	h.UniRoom.Clients.RLock()
	defer h.UniRoom.Clients.RUnlock()
	for userID, client := range h.UniRoom.Clients.items {
		if userID != "" {
			usersID[userID] = client
		}
	}

	return usersID
}

func (h *Hub) IsUserOnline(userID string) bool {
	h.rooms.RLock()
	defer h.rooms.RUnlock()
	_, ok := h.UniRoom.Clients.Get(userID)
	return ok
}

func (h *Hub) GetUsersClient(userID string) (*Client, bool) { // todo get rid of this
	h.rooms.RLock()
	defer h.rooms.RUnlock()
	h.UniRoom.Clients.RLock()
	defer h.UniRoom.Clients.RUnlock()
	for clientUserID, client := range h.UniRoom.Clients.items {
		if clientUserID == userID {
			return client, true
		}
	}
	return nil, false
}

func (h *Hub) isThereRoom(room *Room) bool {
	_, ok := h.rooms.items[room.ID]
	return ok
}

func (h *Hub) GetRoom(id string) (*Room, bool) {
	h.rooms.RLock()
	defer h.rooms.RUnlock()
	room, ok := h.rooms.items[id]
	return room, ok
}

func (h *Hub) BroadcastMessageInRoom(content json.RawMessage, room *Room)  map[string]bool{
	message := &message{
		content:   content,
		room:      room,
		usersSentTo: make(chan map[string]bool),
	}

	h.broadcast <- message
	return <-message.usersSentTo	
}

func (h *Hub) SendMessageToUsers(content json.RawMessage, usersID []string)  map[string]bool {
	message := &messageToListOfUsers{
		content: content,
		usersID: usersID, // TODO change to chan as above
		usersSentTo: make(chan map[string]bool),
	}

	h.sendToListOfUsers <- message
	return <-message.usersSentTo	
}
