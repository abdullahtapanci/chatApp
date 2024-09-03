package main

import (
	"exclusiveChat/api"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	// CORS middleware configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin,Authorization , Content-Type, Accept, Group, User",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}))

	app.Static("/userImages", "./userImages")

	api.SetupRoutes(app)

	go api.SocketHandler()

	err := app.Listen(":8000")
	if err != nil {
		panic(err)
	}

}
