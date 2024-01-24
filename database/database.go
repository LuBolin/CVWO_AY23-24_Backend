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

	fmt.Println("DSN: ", dsn)

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

// MySQL legacy code
// import (
// 	"database/sql"
// 	"fmt"
// 	"os"

// 	_ "github.com/go-sql-driver/mysql"
// )

// var conn *sql.DB // Caps first letter to export variable

// func InitDB() *sql.DB {
// 	db_user := os.Getenv("DB_USER")
// 	db_pwd := os.Getenv("DB_PASSWORD")
// 	db_host := os.Getenv("DB_HOST")
// 	db_port := os.Getenv("DB_PORT")
// 	db_name := os.Getenv("DB_NAME")
// 	dsn := db_user + ":" + db_pwd +
// 		"@tcp(" + db_host + ":" + db_port + ")/" + db_name

// 	fmt.Println("Dsn: ", dsn)
// 	var err error
// 	conn, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		panic("Failed to connect to database")
// 	}

// 	err = conn.Ping()
// 	if err != nil {
// 		panic("Failed to ping database with error" + err.Error())
// 	}

// 	conn.SetMaxOpenConns(10)
// 	conn.SetMaxIdleConns(10)

// 	return conn
// }
