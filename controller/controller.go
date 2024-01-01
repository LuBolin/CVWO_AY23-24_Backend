package controller

import (
	"cvwo/hash"
	"cvwo/jwt"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// username, password(unhashed), email
func SignUp(c *gin.Context, db_conn *sql.DB) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password:" + password + "\n")

	var usernameExists, emailExists int
	query := `SELECT
		MAX(CASE WHEN username = ? THEN 1 ELSE 0 END) AS username_exists,
		MAX(CASE WHEN email = ? THEN 1 ELSE 0 END) AS email_exists
  		FROM Users;`
	db_conn.QueryRow(query, username, email).
		Scan(&usernameExists, &emailExists) // assigns query values to variables

	if usernameExists == 1 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "Username already exists"})
		return
	}
	if emailExists == 1 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "Email already exists"})
		return
	}

	hashed, salt := hash.GenHash(c, password)

	// transaction
	tx, err := db_conn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err})
		return
	}

	var sqlError = func() {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err})
	}

	userInsertQuery := "INSERT INTO Users (username, email) VALUES (?, ?)"
	result, err := tx.Exec(userInsertQuery, username, email)
	if err != nil {
		sqlError()
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		sqlError()
		return
	}

	passwordInsertQuery := "INSERT INTO Passwords (user_id, hashed, salt) VALUES (?, ?, ?)"
	_, err = tx.Exec(passwordInsertQuery, userID, hashed, salt)
	if err != nil {
		sqlError()
		return
	}

	err = tx.Commit()
	if err != nil {
		sqlError()
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success"})
}

func SignIn(c *gin.Context, db_conn *sql.DB) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var hashed []byte
	var salt []byte

	query := `SELECT p.hashed, p.salt
    	FROM Passwords p
    	JOIN Users u ON p.user_id = u.id
    	WHERE u.username = ?
	`
	err := db_conn.QueryRow(query, username).Scan(&hashed, &salt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Username does not exist"})
		return
	} else if err != nil {
		// fmt.Printf("SQL Error: %v\n", err) // %v: value in the default format
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Internal Server Error"})
		return
	}

	match := hash.CompareHash(c, password, hashed, salt)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Wrong password"})
		return
	}

	tokenString := jwt.GenToken(c, username)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "jwt_auth": tokenString})
}

func GetHome(c *gin.Context, db_conn *sql.DB) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Home Success"})
}

func GetMapView(c *gin.Context, db_conn *sql.DB) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "MapView Success"})
}
