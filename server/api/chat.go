package api

import (
	"exclusiveChat/database"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ***************Chat***************
func GetFriends(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

	userId := claims["id"].(string)
	userIdInt, _ := strconv.Atoi(userId)

	//open db connection
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	query := `
		SELECT u.id, u.firstName, u.lastName, u.userName, u.profileImageName
		FROM exclusiveChatDB.users u
		INNER JOIN exclusiveChatDB.friendships f ON u.id = f.user2_id
		WHERE f.user1_id = ?
	`

	// Execute the SQL query
	rows, err := db.Query(query, userIdInt)
	if err != nil {
		fmt.Println("Error executing query:", err)
	}
	defer rows.Close()

	// Collect the fetched rows into a slice of structs
	var friends []UserData
	for rows.Next() {
		var friend UserData
		if err := rows.Scan(&friend.Id, &friend.FirstName, &friend.LastName, &friend.UserName, &friend.ProfileImageName); err != nil {
			fmt.Println("Error scanning row:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Error scanning row"})
		}
		friends = append(friends, friend)
	}

	return c.JSON(fiber.Map{
		"message": "process successful",
		"friends": friends,
	})
}
