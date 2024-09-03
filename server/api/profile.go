package api

import (
	"exclusiveChat/database"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ***************PROFİLE***************
func Profile(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

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
		fmt.Println("error retriving user data lskmskmdx")
		fmt.Println(err)
	}

	// Return data
	return c.JSON(fiber.Map{
		"message":          "User datas fetched successfuly",
		"firstName":        fetchedUser.FirstName,
		"lastName":         fetchedUser.LastName,
		"userName":         fetchedUser.UserName,
		"email":            fetchedUser.Email,
		"profileImageName": fetchedUser.ProfileImageName,
	})
}

// ***************Update Profile Info***************
func UpdateProfileInfo(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

	userId := claims["id"].(string)

	var updatedData UserData

	if err := c.BodyParser(&updatedData); err != nil {
		fmt.Println("Invalid JSON")
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	//open database and close when precess is done
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	// Prepare the UPDATE statement
	stmt, err := db.Prepare("UPDATE exclusiveChatDB.users SET firstName = ?, lastName = ?, email = ?, userName = ? WHERE id = ?")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while preparing the update statement"})
	}
	defer stmt.Close()

	// Execute the UPDATE statement with updated data
	_, err = stmt.Exec(updatedData.FirstName, updatedData.LastName, updatedData.Email, updatedData.UserName, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while executing the update statement"})
	}

	//generate jwt token
	tokenString, err := generateJWTToken(userId, updatedData.UserName)
	if err != nil {
		fmt.Println("Failed to generte JWT token")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate JWT token")
	}

	// Return data
	return c.JSON(fiber.Map{
		"message":   "User datas updated successfuly",
		"firstName": updatedData.FirstName,
		"lastName":  updatedData.LastName,
		"userName":  updatedData.UserName,
		"email":     updatedData.Email,
		"token":     tokenString,
	})
}

// ***************Update Profile Image***************
func UpdateProfileImage(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

	userId := claims["id"].(string)

	// Parse the form data containing the file
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("error parsing the form", err)
		return err
	}

	// Get the file from form data
	files := form.File["image"]
	if len(files) == 0 {
		return fmt.Errorf("no files uploaded")
	}

	// Get the first file (assuming single file upload)
	file := files[0]

	// Generate a unique filename for the image
	newUserProfileImageName := uuid.New().String() + ".png"

	// Save the file to the server
	if err := c.SaveFile(file, "./userImages/"+newUserProfileImageName); err != nil {
		fmt.Println("error saving the image", err)
		return err
	}

	//open db connection
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	var fetchedUser UserData

	// Execute the query to retrieve specific data
	err = db.QueryRow("SELECT profileImageName FROM exclusiveChatDB.users WHERE id = ?", userId).Scan(&fetchedUser.ProfileImageName)

	if err != nil {
		fmt.Println("error retriving user data")
		fmt.Println(err)
	}

	if len(fetchedUser.ProfileImageName) > 0 {
		// Construct the file path
		filePath := fmt.Sprintf("./userImages/%s", fetchedUser.ProfileImageName)

		// Attempt to delete the file
		if err := os.Remove(filePath); err != nil {
			// If an error occurs while deleting the file, return an error response
			fmt.Println("failed to delete image", err)
		}

		fmt.Println("former image deleted")
	}

	// Prepare the UPDATE statement
	stmt, err := db.Prepare("UPDATE exclusiveChatDB.users SET profileImageName = ? WHERE id = ?")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while preparing the update statement"})
	}
	defer stmt.Close()

	// Execute the UPDATE statement with updated data
	_, err = stmt.Exec(newUserProfileImageName, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while executing the update statement"})
	}

	// Return data
	return c.JSON(fiber.Map{
		"message":          "User profile image saved successfuly",
		"profileImageName": newUserProfileImageName,
	})

}

func GetUserName(c *fiber.Ctx) error {
	_, claims, err := ValidateToken(c)
	if err != nil {
		fmt.Println(err) // Return the error if token validation fails
	}

	userId := claims["id"].(string)

	var data UserData

	//open database and close when precess is done
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
		return c.Status(400).JSON(fiber.Map{"error": "Error opening database connection:"})
	}

	defer database.CloseDBConnection(db)

	err = db.QueryRow("SELECT id, userName FROM exclusiveChatDB.users WHERE id = ?", userId).Scan(&data.Id, &data.UserName)

	if err != nil {
		fmt.Println("error retriving user data")
		fmt.Println(err)
	}

	// Return data
	return c.JSON(fiber.Map{
		"userName": data.UserName,
		"id":       data.Id,
	})
}
