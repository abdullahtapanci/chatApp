package api

import (
	"fmt"
	"log"
	"strconv"

	"exclusiveChat/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SocketHandler() {
	for {
		select {
		case client := <-ws.Register:
			// Pre-initialize organization map if it doesn't exist
			if ws.Clients[client.Group] == nil {
				ws.Clients[client.Group] = make(map[string]*websocket.Conn)
			}
			ws.Clients[client.Group][client.User] = client.Conn
			log.Println("client registered:", client.Group, client.User)

		case message := <-ws.Broadcast:
			for org, users := range ws.Clients {
				if org == message.From.Group {
					if len(users) <= 1 {
						//user is not online
						//save message to send later
						AddMessage(ws.MessageQueue, message.To, message.From.UserId, message.Msg)
					}
					for user, conn := range users {
						if org != message.From.Group || user != message.From.User {
							if err := conn.WriteMessage(websocket.TextMessage, []byte(message.Msg)); err != nil {
								log.Println("write error:", err)
								ws.RemoveClient(org, user) // Update client removal
								conn.WriteMessage(websocket.CloseMessage, []byte{})
								conn.Close()
							}
						}
					}
				}
			}

		case client := <-ws.Unregister:
			ws.RemoveClient(client.Group, client.User) // Update client removal
			log.Println("client unregistered:", client.Group, client.User)
		}
	}
}

func SendQueuedMessages(c *fiber.Ctx) error {

	group := c.Query("Group")
	user := c.Query("User")
	friendId, err := strconv.Atoi(c.Query("SelectedFriendId"))
	if err != nil {
		fmt.Println("Error converting SelectedFriendId to int:", err)
	}
	userId, err := strconv.Atoi(c.Query("UserId"))
	if err != nil {
		fmt.Println("Error converting UserId to int:", err)
	}

	conn := ws.Clients[group][user]
	if conn == nil {
		log.Println("WebSocket connection is nil")
	}

	// Collect keys to delete
	var keysToDelete []int

	for id, messages := range ws.MessageQueue {
		if id == userId {
			value, exists := messages[friendId]
			if exists {
				for _, message := range value {
					// Send message to the WebSocket connection
					err := conn.WriteMessage(websocket.TextMessage, []byte(message))
					if err != nil {
						log.Println("Error sending message:", err)
						// Handle error (optional: close the connection or notify the user)
					}
				}
				// Mark the key for deletion
				keysToDelete = append(keysToDelete, friendId)
			}
		}
	}

	// Delete the messages outside of the iteration
	for _, id := range keysToDelete {
		delete(ws.MessageQueue[userId], id)
	}
	if len(ws.MessageQueue[userId]) == 0 {
		delete(ws.MessageQueue, userId)
	}
	return c.JSON(fiber.Map{
		"message": "successfull",
	})

}

func CheckIfThereAreNewMessages(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

	userId, err := strconv.Atoi(claims["id"].(string))
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}

	type newMessagesMap map[int]int //friend id , message numeber
	newMessages := make(newMessagesMap)

	if len(ws.MessageQueue[userId]) != 0 {
		for friendId, messages := range ws.MessageQueue[userId] {
			newMessages[friendId] = len(messages)
		}
	}

	fmt.Println(newMessages)

	return c.JSON(fiber.Map{
		"data": newMessages,
	})

}

func AddMessage(queue ws.MessageQueueStruct, userId int, friendId int, message string) {
	// Check if the outer map contains the key for the user
	if _, exists := queue[userId]; !exists {
		// Initialize the outer map entry if it does not exist
		queue[userId] = make(map[int][]string)
	}

	// Check if the inner map contains the key for the friend
	if _, exists := queue[userId][friendId]; !exists {
		// Initialize the inner map entry if it does not exist
		queue[userId][friendId] = []string{}
	}

	// Append the new message to the slice for the specific friend
	queue[userId][friendId] = append(queue[userId][friendId], message)
}
