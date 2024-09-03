package api

import (
	"exclusiveChat/middleware"
	"exclusiveChat/ws"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Create user data struct
type UserData struct {
	Id               string `json:"id"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	UserName         string `json:"userName"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	ProfileImageName string `json:"profileImageName"`
}

// create func that set routes
func SetupRoutes(app *fiber.App) {
	app.Post("/api/signup", SignUp)
	app.Post("/api/login", Login)
	app.Post("/api/home", Home)
	app.Post("/api/profile", Profile)
	app.Post("/api/updateProfileInfo", UpdateProfileInfo)
	app.Post("/api/updateProfileImage", UpdateProfileImage)
	app.Post("/api/findUser", FindUser)
	app.Post("/api/createFriendship", CreateFriendship)
	app.Post("/api/getFriends", GetFriends)
	app.Get("/api/getUserName", GetUserName)
	app.Get("/api/CheckMessages", CheckIfThereAreNewMessages)
	app.Get("/api/getNewMessages", SendQueuedMessages)
	app.Get("/ws", middleware.WsMiddleware, websocket.New(func(c *websocket.Conn) {

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				c.WriteJSON(fiber.Map{"customError": "error occurred"})
			}
			c.Close()
		}()

		userId, err := strconv.Atoi(c.Locals("UserId").(string))
		if err != nil {
			log.Println("Error converting userId to int:", err)
			c.WriteJSON(fiber.Map{"customError": "invalid userId"})
		}

		clientObj := ws.ClientObject{
			Group:  c.Locals("Group").(string),
			User:   c.Locals("User").(string),
			UserId: userId,
			Conn:   c,
		}
		defer func() {
			ws.Unregister <- clientObj
			c.Close()
		}()
		// Register the client
		ws.Register <- clientObj

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				to, err := strconv.Atoi(c.Locals("To").(string))
				if err != nil {
					log.Println("Error converting To to int:", err)
					c.WriteJSON(fiber.Map{"customError": "invalid To"})
					continue
				}
				// Broadcast the received message
				ws.Broadcast <- ws.BroadcastObject{
					Msg:  string(message),
					From: clientObj,
					To:   to,
				}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))
}
