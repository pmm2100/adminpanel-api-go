package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	godotenv.Load()
	var err error
	DB, err = sql.Open("sqlite", os.Getenv("DB_PATH"))
	if err != nil {
		panic(fmt.Errorf("Could not initiate database!\n %v", err))
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()

}

func createTables() {

	createTableUsers := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			isadmin INTEGER NOT NULL DEFAULT 0,
			isbanned INTEGER NOT NULL DEFAULT 0
		)
	`

	_, err := DB.Exec(createTableUsers)

	if err != nil {
		panic(fmt.Errorf("Could not create users table! \n %v", err))
	}

	createTablePermBanned := `
		CREATE TABLE IF NOT EXISTS permBanned (
			id INTEGER PRIMARY KEY,
			userIP TEXT NOT NULL UNIQUE
		)
	`

	_, err = DB.Exec(createTablePermBanned)

	if err != nil {
		panic(fmt.Errorf("Could not create permBanned table! \n %v", err))
	}

}
