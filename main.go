package main

import (
	"cvwo/database"
	"cvwo/router"
	"database/sql"
	"log"

	"github.com/joho/godotenv"
)

var db_conn *sql.DB

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	db_conn = database.InitDB()
	router := router.InitRouter(db_conn)

	router.Run(":8080")
}
