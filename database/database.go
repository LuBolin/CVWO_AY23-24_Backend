package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var conn *sql.DB // Caps first letter to export variable

func InitDB() *sql.DB {
	// data source name
	db_user := "root"
	db_pwd := "password"
	db_host := "localhost"
	db_port := "3306"
	db_name := "CVWO"
	dsn := db_user + ":" + db_pwd +
		"@tcp(" + db_host + ":" + db_port + ")/" + db_name

	var err error
	conn, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect to database")
	}

	err = conn.Ping()
	if err != nil {
		panic("Failed to ping database")
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	return conn
}
