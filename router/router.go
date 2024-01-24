package router

import (
	"cvwo/controller"
	"cvwo/jwt"
	"database/sql"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var db_conn *sql.DB

func InitRouter(incoming_db_conn *sql.DB) *gin.Engine {
	db_conn = incoming_db_conn
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_IP")}, // Front end server ip
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	accGrp := router.Group("/account")
	{
		accGrp.POST("/signin", func(c *gin.Context) { controller.SignIn(c, db_conn) })
		accGrp.POST("/signup", func(c *gin.Context) { controller.SignUp(c, db_conn) })
	}

	forumRoutes := router.Group("")
	{
		forumRoutes.GET("/forum", func(c *gin.Context) { controller.GetForum(c, db_conn) })
		forumRoutes.GET("/post/:post_id", func(c *gin.Context) { controller.GetPost(c, db_conn) })
	}

	authRoute := router.Group("/auth").Use(jwt.ValidateJWT)
	{
		authRoute.POST("/check", func(c *gin.Context) {
			controller.CheckJWT(c, db_conn)
		})
		authRoute.POST("/newpost", func(c *gin.Context) {
			controller.NewPost(c, db_conn)
		})
		authRoute.POST("/comment", func(c *gin.Context) {
			controller.NewComment(c, db_conn)
		})
		authRoute.POST("/updatecomment/", func(c *gin.Context) {
			controller.UpdateComment(c, db_conn)
		})
		authRoute.POST("/deletecomment/", func(c *gin.Context) {
			controller.DeleteComment(c, db_conn)
		})
	}

	return router
}
