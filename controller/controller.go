package controller

import (
	"cvwo/hash"
	"cvwo/jwt"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PostForm(x)'s x has to match front ebd's URLSearchParams's variable name

func SignUp(c *gin.Context, db_conn *sql.DB) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password:" + password + "\n")

	var usernameExists, emailExists int
	query := `
		SELECT
			CASE WHEN EXISTS (SELECT 1 FROM Users WHERE username = $1) THEN 1 ELSE 0 END AS username_exists,
			CASE WHEN EXISTS (SELECT 1 FROM Users WHERE email = $2) THEN 1 ELSE 0 END AS email_exists
	`
	err := db_conn.QueryRow(query, username, email).Scan(&usernameExists, &emailExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	if usernameExists == 1 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "Username already exists"})
		return
	}
	if emailExists == 1 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "Email already exists"})
		return
	}

	hashed, salt := hash.GenHash(c, password)

	tx, err := db_conn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	var sqlError = func() {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Internal Server Error"})
	}

	userInsertQuery := "INSERT INTO Users (username, email) VALUES ($1, $2) RETURNING id"
	row := tx.QueryRow(userInsertQuery, username, email)
	var userID int
	if err := row.Scan(&userID); err != nil {
		sqlError()
		return
	}

	passwordInsertQuery := "INSERT INTO Passwords (user_id, hashed, salt) VALUES ($1, $2, $3)"
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

	var user_id int
	var hashed []byte
	var salt []byte

	query := `
		SELECT u.id, p.hashed, p.salt
		FROM Passwords p
		JOIN Users u ON p.user_id = u.id
		WHERE u.username = $1
	`
	err := db_conn.QueryRow(query, username).Scan(&user_id, &hashed, &salt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Username does not exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Internal Server Error"})
		return
	}

	match := hash.CompareHash(c, password, hashed, salt)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Wrong password"})
		return
	}

	tokenString := jwt.GenToken(c, username, user_id)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "jwtToken": tokenString})
}

type PostData struct {
	PostID  int    `json:"post_id"`
	Title   string `json:"title"`
	Topic   string `json:"topic"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

type CommentData struct {
	CommentID  int    `json:"comment_id"`
	AuthorName string `json:"author"`
	AuthorID   int    `json:"author_id"`
	Date       string `json:"date"`
	Content    string `json:"content"`
}

func GetForum(c *gin.Context, db_conn *sql.DB) {
	title := c.DefaultQuery("title", "")
	topic := c.DefaultQuery("topic", "All")
	offsetStr := c.DefaultQuery("offset", "0")
	chunkSizeStr := c.DefaultQuery("chunksize", "0")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chunksize parameter"})
		return
	}

	baseQuery :=
		`SELECT Posts.id, Posts.title, Posts.topic, Users.username, 
			TO_CHAR(Posts.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM Posts
		JOIN Users ON Posts.author_id = Users.id
		WHERE 1=1`

	var args []interface{}
	var conditions []string
	argNumber := 1

	if title != "" {
		conditions = append(conditions, fmt.Sprintf("Posts.title ILIKE $%d", argNumber))
		args = append(args, "%"+title+"%")
		argNumber++
	}

	if topic != "All" {
		conditions = append(conditions, fmt.Sprintf("Posts.topic = $%d", argNumber))
		args = append(args, topic)
		argNumber++
	}

	queryCondition := ""
	if len(conditions) > 0 {
		queryCondition = " AND " + strings.Join(conditions, " AND ")
	}

	limitArgNumber := argNumber
	offsetArgNumber := argNumber + 1

	finalQuery := fmt.Sprintf("%s%s LIMIT $%d OFFSET $%d", baseQuery, queryCondition, limitArgNumber, offsetArgNumber)
	args = append(args, chunkSize, offset*chunkSize)

	rows, err := db_conn.Query(finalQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
		return
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var (
			id         int
			title      string
			topic      string
			author     string
			created_at string
		)
		err := rows.Scan(&id, &title, &topic, &author, &created_at)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		post := PostData{
			PostID: id,
			Title:  title,
			Topic:  topic,
			Author: author,
			Date:   created_at,
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating the rows"})
		return
	}

	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM Posts JOIN Users ON Posts.author_id = Users.id WHERE 1=1%s", queryCondition)

	var totalPosts int
	err = db_conn.QueryRow(countQuery, args[0:len(args)-2]...).Scan(&totalPosts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying total post count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts, "post_count": totalPosts})
}

func NewPost(c *gin.Context, db_conn *sql.DB) {
	user_id := jwt.GetUserIdFromToken(c)
	title := c.PostForm("title")
	topic := c.PostForm("topic")
	content := c.PostForm("content")
	postQuery := `
        INSERT INTO Posts (author_id, title, topic, content)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var postID int
	err := db_conn.QueryRow(postQuery, user_id, title, topic, content).Scan(&postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Error creating the post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "new_post_id": postID})
}

func GetPost(c *gin.Context, db_conn *sql.DB) {
	postID := c.DefaultQuery("post_id", "-1")
	offsetStr := c.DefaultQuery("offset", "0")
	chunkSizeStr := c.DefaultQuery("chunksize", "0")

	if postID == "-1" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chunksize parameter"})
		return
	}

	postQuery := `
        SELECT 
            p.title, p.topic, u.username, p.content,
            TO_CHAR(p.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
        FROM Posts p
        JOIN Users u ON p.author_id = u.id
        WHERE p.id = $1
    `
	row := db_conn.QueryRow(postQuery, postID)
	var (
		title      string
		topic      string
		author     string
		content    string
		created_at string
	)
	err = row.Scan(&title, &topic, &author, &content, &created_at)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Error getting the post"})
		return
	}

	postIDInt, _ := strconv.Atoi(postID)
	post := PostData{
		PostID:  postIDInt,
		Title:   title,
		Topic:   topic,
		Author:  author,
		Content: content,
		Date:    created_at,
	}

	commentsQuery := `
        SELECT 
            c.id, c.content, u.username, u.id,
            TO_CHAR(c.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
        FROM Comments c
        JOIN Users u ON c.author_id = u.id
        WHERE c.post_id = $1
		ORDER BY created_at ASC
        LIMIT $2 OFFSET $3
    `

	limit := chunkSize
	offset = offset * chunkSize
	rows, err := db_conn.Query(commentsQuery, postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying comments"})
		return
	}
	defer rows.Close()

	var comments []CommentData
	for rows.Next() {
		var (
			commentID  int
			authorName string
			authorID   int
			date       string
			content    string
		)
		err := rows.Scan(&commentID, &content, &authorName, &authorID, &date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning comments"})
			return
		}

		comment := CommentData{
			CommentID:  commentID,
			AuthorName: authorName,
			AuthorID:   authorID,
			Date:       date,
			Content:    content,
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating comments"})
		return
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM Comments 
		WHERE post_id = $1
	`

	var totalComments int
	err = db_conn.QueryRow(countQuery, postID).Scan(&totalComments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying total post count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post, "comments": comments, "comment_count": totalComments})
}

func NewComment(c *gin.Context, db_conn *sql.DB) {
	user_id := jwt.GetUserIdFromToken(c)
	postID := c.PostForm("post_id")
	content := c.PostForm("content")
	commentQuery := `
		INSERT INTO Comments (author_id, post_id, content)
		VALUES ($1, $2, $3)
	`
	_, err := db_conn.Exec(commentQuery, user_id, postID, content)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Error creating the comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200})
}

func UpdateComment(c *gin.Context, db_conn *sql.DB) {
	user_id := jwt.GetUserIdFromToken(c)
	postID := c.PostForm("post_id")
	commentID := c.PostForm("comment_id")
	content := c.PostForm("content")

	// Check if the comment belongs to the user trying to update it
	var dbUserID int
	checkOwnerQuery := "SELECT author_id FROM Comments WHERE id = $1 AND post_id = $2"
	err := db_conn.QueryRow(checkOwnerQuery, commentID, postID).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Comment not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Error checking comment owner"})
		return
	}

	if dbUserID != user_id {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Unauthorized to update this comment"})
		return
	}

	updateQuery := `
		UPDATE Comments 
		SET content = $1
		WHERE id = $2 AND post_id = $3
	`
	_, err = db_conn.Exec(updateQuery, content, commentID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Error updating the comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Comment updated successfully"})
}

func DeleteComment(c *gin.Context, db_conn *sql.DB) {
	user_id := jwt.GetUserIdFromToken(c)

	c.Request.ParseForm() // couldnt get DELETE to work, sad
	postID := c.PostForm("post_id")
	commentID := c.PostForm("comment_id")

	// Check if the comment belongs to the user trying to delete it
	var dbUserID int
	checkOwnerQuery := "SELECT author_id FROM Comments WHERE id = $1 AND post_id = $2"
	err := db_conn.QueryRow(checkOwnerQuery, commentID, postID).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Comment not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Error checking comment owner"})
		return
	}

	if dbUserID != user_id {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Unauthorized to delete this comment"})
		return
	}

	deleteQuery := `
		DELETE FROM Comments
		WHERE id = $1 AND post_id = $2
	`
	_, err = db_conn.Exec(deleteQuery, commentID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Error deleting the comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Comment deleted successfully"})
}

func CheckJWT(c *gin.Context, db_conn *sql.DB) {
	// already authenticated by ValidateJWT middleware
	user_id := jwt.GetUserIdFromToken(c)
	c.JSON(http.StatusOK, gin.H{"user_id": user_id})
}
