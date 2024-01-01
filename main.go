package main

import (
	"cvwo/database"
	"cvwo/router"
	"database/sql"
)

var db_conn *sql.DB

func main() {
	db_conn = database.InitDB()
	router := router.InitRouter(db_conn)

	router.Run(":8080")
}
