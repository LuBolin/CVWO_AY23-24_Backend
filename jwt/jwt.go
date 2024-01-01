package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("111")

func GenToken(c *gin.Context, username string) (tokenString string) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ = token.SignedString(hmacSampleSecret)

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
		return hmacSampleSecret, nil
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
