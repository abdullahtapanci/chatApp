package api

import (
	"exclusiveChat/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ***************HOME***************
func Home(c *fiber.Ctx) error {
	// Validate the token
	token, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
		return err
	}

	// Token is valid, you can access claims and proceed with your logic
	// Example: Access the user's name from the token claims
	userId := claims["id"].(string)

	var fetchedUser UserData

	//open database and close when precess is done
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	// Execute the query to retrieve specific data
	err = db.QueryRow("SELECT id, firstName, lastName, userName, email, profileImageName FROM exclusiveChatDB.users WHERE id = ?", userId).Scan(&fetchedUser.Id, &fetchedUser.FirstName, &fetchedUser.LastName, &fetchedUser.UserName, &fetchedUser.Email, &fetchedUser.ProfileImageName)

	if err != nil {
		fmt.Println("error retriving user data")
		fmt.Println(err)
	}

	// Return success response
	return c.JSON(fiber.Map{
		"message":          "Token validated successfully",
		"firstName":        fetchedUser.FirstName,
		"lastName":         fetchedUser.LastName,
		"userName":         fetchedUser.UserName,
		"profileImageName": fetchedUser.ProfileImageName,
		"parsedToken":      token,
	})

}
