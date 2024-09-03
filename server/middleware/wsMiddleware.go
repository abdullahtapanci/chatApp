package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WsMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		// Your authentication process goes here. Get the Token from header and validate it
		// Extract the claims from the token and set them to the Locals
		// This is because you cannot access headers in the websocket.Conn object below
		group := c.Query("Group")
		user := c.Query("User")
		to := c.Query("SelectedFriendId")
		userId := c.Query("UserId")

		c.Locals("Group", group)
		c.Locals("User", user)
		c.Locals("To", to)
		c.Locals("UserId", userId)

		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
