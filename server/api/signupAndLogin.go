package api

import (
	"database/sql"
	"exclusiveChat/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ***************SİGN UP***************
func SignUp(c *fiber.Ctx) error {
	// Parse JSON request body into user data struct
	var userData UserData
	if err := c.BodyParser(&userData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	//open database and close when precess is done
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	//add user to database
	_, err := db.Query("INSERT INTO `exclusiveChatDB`.`users` (`firstName`,`lastName`, `userName`,`email`,`password`,`profileImageName`) VALUES ( ?, ?, ?, ?, ?, ?)",
		userData.FirstName, userData.LastName, userData.UserName, userData.Email, userData.Password, "null")

	if err != nil {
		fmt.Println("Error inserting user into the database:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Error inserting user into the database"})
	}

	// Send a response
	return c.JSON(fiber.Map{"message": "Signup successful"})
}

// ***************LOG IN***************
func Login(c *fiber.Ctx) error {
	// Parse JSON request body into UserData struct
	var userData UserData
	if err := c.BodyParser(&userData); err != nil {
		fmt.Println("Invalid JSON")
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	//get user data from database if it is exist and compare

	//open database and close when precess is done
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	//create query
	query := `SELECT id, firstName, lastName, userName, email, password FROM exclusiveChatDb.users WHERE email = ?`

	// Execute the query
	var user UserData
	err := db.QueryRow(query, userData.Email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No user found with that username")
		return c.Status(400).JSON(fiber.Map{"error": "No user found"})
	case err != nil:
		fmt.Println(err)
	}

	//check wheter informations match
	if user.Email == userData.Email && user.Password == userData.Password {
		//generate jwt token
		tokenString, err := generateJWTToken(user.Id, user.UserName)
		if err != nil {
			fmt.Println("Failed to generte JWT token")
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate JWT token")
		}

		// Respond with success message and token
		return c.JSON(fiber.Map{"message": "Login successful", "token": tokenString})
	} else {
		return c.JSON(fiber.Map{"message": "Signup unsuccessful. Email or password is incorrect"})
	}

}
