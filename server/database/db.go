package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDBConnection() (*sql.DB, error) {

	dbUsers, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/dbName")

	if err != nil {
		fmt.Println("Error opening the database:", err)
		return nil, err
	}

	return dbUsers, nil
}

func CloseDBConnection(db *sql.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing the database:", err)
		}
	}
}
