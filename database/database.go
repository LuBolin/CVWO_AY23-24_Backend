package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var conn *sql.DB

func InitDB() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// PostgreSQL connection string
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPwd, dbHost, dbPort, dbName,
	)

	var err error
	conn, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	err = conn.Ping()
	if err != nil {
		panic("Failed to ping the database with error: " + err.Error())
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	return conn
}
