package ws

import "github.com/gofiber/websocket/v2"

type MiniClient map[string]map[string]*websocket.Conn // Modified type

type ClientObject struct {
	Group  string
	User   string
	UserId int
	Conn   *websocket.Conn
}

type BroadcastObject struct {
	Msg  string
	From ClientObject
	To   int
}

type MessageQueueStruct map[int]map[int][]string // user id   friend id  messages

var MessageQueue = make(MessageQueueStruct)

var Clients = make(MiniClient) // Initialized as a nested map
var Register = make(chan ClientObject)
var Broadcast = make(chan BroadcastObject)
var Unregister = make(chan ClientObject)

func RemoveClient(org string, user string) {
	if conn, ok := Clients[org][user]; ok { // Check if client exists
		delete(Clients[org], user)
		conn.Close() // Close the connection before potentially removing the organization map
		if len(Clients[org]) == 0 {
			delete(Clients, org) // Remove empty organization map
		}
	}
}
