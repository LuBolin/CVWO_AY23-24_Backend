package router

import (
	"cvwo/controller"
	"cvwo/jwt"
	"database/sql"
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
		AllowOrigins:     []string{"http://localhost:5173"}, // Replace with your client's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	accGrp := router.Group("/account")
	{
		accGrp.POST("/signup", func(c *gin.Context) { controller.SignUp(c, db_conn) })
		accGrp.POST("/signin", func(c *gin.Context) { controller.SignIn(c, db_conn) })
	}

	homeGrp := router.Group("/home").Use(jwt.ValidateJWT)
	{
		homeGrp.GET("", func(c *gin.Context) { controller.GetHome(c, db_conn) })
		homeGrp.GET("/mapview", func(c *gin.Context) { controller.GetMapView(c, db_conn) })
	}

	return router
}
