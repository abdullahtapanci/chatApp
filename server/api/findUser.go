package api

import (
	"encoding/json"
	"exclusiveChat/database"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ***************Find User***************
func FindUser(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

	userId := claims["id"].(string)
	userIdInt, _ := strconv.Atoi(userId)

	var strToBeSearch string

	if err := c.BodyParser(&strToBeSearch); err != nil {
		fmt.Println("Invalid JSON")
	}

	var usersFound []foundUserStruct

	if strToBeSearch == "" || strToBeSearch == " " || strToBeSearch == "  " {
	} else {
		usersFound = searchQuery(strToBeSearch, userIdInt)
	}

	// Marshal the array of users into JSON
	jsonData, err := json.Marshal(usersFound)
	if err != nil {
		fmt.Println("error converting array of users to json", err)
	}

	// Return data
	return c.JSON(fiber.Map{
		"message":    "User profile image saved successfuly",
		"foundUsers": string(jsonData),
	})
}

// ***************Create Friendship***************
func CreateFriendship(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

	userId := claims["id"].(string)
	userIdInt, _ := strconv.Atoi(userId)

	type FriendIdData struct {
		FriendId string `json:"friendId"`
	}

	var friendIdData FriendIdData

	if err := c.BodyParser(&friendIdData); err != nil {
		fmt.Println("Invalid JSON", err)
	}

	friendIdInt, _ := strconv.Atoi(friendIdData.FriendId)

	//open db connection
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	//create a friendship
	_, err = db.Query("INSERT INTO exclusiveChatDB.friendships (user1_id, user2_id) VALUES (?, ?)", userIdInt, friendIdInt)

	if err != nil {
		fmt.Println("Error creating friendship in the database:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Error inserting user into the database"})
	}

	return c.JSON(fiber.Map{
		"message": "became friend successfuly",
	})

}
