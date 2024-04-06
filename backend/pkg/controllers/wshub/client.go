package wshub

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	UserID   string
	UserName string
	Room     *Room

	// The websocket connection.
	Conn *websocket.Conn
	// Buffered channel of received messages.
	ReceivedMessages chan []byte
	// TODO chan for hub errors

	Registered chan bool
}

// TODO use this in the correct place
func NewClient(hub *Hub, userID, userName string, room *Room, conn *websocket.Conn, receivedMessages chan []byte, clientRegistered chan bool)( *Client, bool) {
	client := &Client{
		UserID:   userID,
		UserName: userName,
		Room:     room,
		Conn:     conn,
	}

	if receivedMessages == nil {
		client.ReceivedMessages = make(chan []byte, 256)
	} else {
		client.ReceivedMessages = receivedMessages
	}

	if clientRegistered == nil {
		client.Registered = make(chan bool)
	} else {
		client.Registered = clientRegistered
	}

	hub.RegisterClientToHub(client)
	// Wait for client registration to complete
	ok:=<-client.Registered
	return client, ok
}

func (c *Client) WriteMessage(message []byte) {

	c.ReceivedMessages <- message
}

func (c *Client) String() string {
	return fmt.Sprintf("addr: %p || User ID:'%s' || connection: %p || channels: clientRegistered %p  |  ReceivedMessages %p", c, c.UserID, c.Conn, c.Registered, c.ReceivedMessages)
}
