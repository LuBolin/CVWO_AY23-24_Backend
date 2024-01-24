package jwt

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// var hmacSampleSecret = []byte("111")
var hmacSecret = []byte(os.Getenv("HMAC_SECRET"))

func GenToken(c *gin.Context, username string, userID int) (tokenString string) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix() // expiry time is 12 hours

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ = token.SignedString(hmacSecret)

	return tokenString
}

func ValidateJWT(c *gin.Context) {
	const BearerSchema = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Auth token required"})
		c.Abort()
		return
	}
	if len(authHeader) <= len(BearerSchema) || authHeader[:len(BearerSchema)] != BearerSchema {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth header format"})
		c.Abort()
		return
	}
	tokenString := authHeader[len(BearerSchema):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth token"})
		c.Abort()
		return
	}

	// claims, ok
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth token"})
		c.Abort()
		return
	}
}

func GetUserIdFromToken(c *gin.Context) int {
	const BearerSchema = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Auth token required"})
		c.Abort()
		return -1
	}
	if len(authHeader) <= len(BearerSchema) || authHeader[:len(BearerSchema)] != BearerSchema {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth header format"})
		c.Abort()
		return -1
	}
	tokenString := authHeader[len(BearerSchema):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth token"})
		c.Abort()
		return -1
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		c.Next()
		user_id := int(claims["user_id"].(float64))
		return user_id
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth token"})
		c.Abort()
		return -1
	}
}

func GetUsernameFromToken(c *gin.Context) int {
	const BearerSchema = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Auth token required"})
		c.Abort()
		return -1
	}
	if len(authHeader) <= len(BearerSchema) || authHeader[:len(BearerSchema)] != BearerSchema {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth header format"})
		c.Abort()
		return -1
	}
	tokenString := authHeader[len(BearerSchema):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth token"})
		c.Abort()
		return -1
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		c.Next()
		user_id := int(claims["username"].(float64))
		return user_id
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid auth token"})
		c.Abort()
		return -1
	}
}
