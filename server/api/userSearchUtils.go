package api

import (
	"database/sql"
	"exclusiveChat/database"
	"fmt"
	"strconv"
)

type foundUserStruct struct {
	FoundUserId        string `json:"id"`
	FoundUserFirstName string `json:"firstName"`
	FoundUserLastName  string `json:"lastName"`
	FoundUserUserName  string `json:"userName"`
	FoundUserImageName string `json:"imageName"`
}

func minForLevenshteinDistance(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= a && b <= c {
		return b
	}
	return c
}

func levenshteinDistance(s1, s2 string) int {
	m := len(s1)
	n := len(s2)

	// Create a 2D slice to store the distances
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Initialize the first row and column
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// Calculate the distances
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + minForLevenshteinDistance(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}

	return dp[m][n]
}

func searchQuery(query string, userId int) []foundUserStruct {
	//open database and close when precess is done
	db, dbErr := database.OpenDBConnection()

	if db == nil {
		fmt.Println("Error opening database connection:", dbErr)
	}

	defer database.CloseDBConnection(db)

	var fetchedUser UserData

	var numberOfUsersInDB int

	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&numberOfUsersInDB)
	if err != nil {
		fmt.Println("error getting number of users stored in db", err)
	}

	var usersFound []foundUserStruct

	var distance1, distance2, distance3 int

	for idToBeFetched := 0; idToBeFetched < numberOfUsersInDB; idToBeFetched++ {
		// Execute the query to retrieve specific data
		err = db.QueryRow("SELECT id, firstName, lastName, userName, profileImageName FROM exclusiveChatDB.users ORDER BY id LIMIT 1 OFFSET ?", idToBeFetched).Scan(&fetchedUser.Id, &fetchedUser.FirstName, &fetchedUser.LastName, &fetchedUser.UserName, &fetchedUser.ProfileImageName)

		if err != nil {
			fmt.Println("error retriving user data")
			fmt.Println(err)
			continue
		}

		num, _ := strconv.Atoi(fetchedUser.Id)

		if num != userId {
			distance1 = levenshteinDistance(fetchedUser.FirstName, query)
			distance2 = levenshteinDistance(fetchedUser.FirstName+" "+fetchedUser.LastName, query)
			distance3 = levenshteinDistance(fetchedUser.UserName, query)

			if distance1 < 5 || distance2 < 5 || distance3 < 5 {
				newUser := foundUserStruct{
					FoundUserId:        fetchedUser.Id,
					FoundUserFirstName: fetchedUser.FirstName,
					FoundUserLastName:  fetchedUser.LastName,
					FoundUserUserName:  fetchedUser.UserName,
					FoundUserImageName: fetchedUser.ProfileImageName,
				}
				usersFound = append(usersFound, newUser)
			}
		}
	}

	return usersFound

}

func CheckUserExistsByID(db *sql.DB, userID int) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE id = ?"

	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
