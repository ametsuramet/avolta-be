package middleware

import (
	"avolta/config"
	"avolta/model"
	"avolta/object/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		splitToken := strings.Split(tokenString, "Bearer ")
		if len(splitToken) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}
		reqToken := splitToken[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token unsplited"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(reqToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.App.Server.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "token unparsed"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			fmt.Println("Authorized user:", claims.Id)
			var user = model.User{Base: model.Base{ID: claims.Id}}
			user.GetUserByID()
			// fmt.Println("user", user)
			c.Set("user", user)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		getUser, _ := c.Get("user")
		user := getUser.(model.User)
		if !user.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "user is not admin"})
			c.Abort()
		}
		user.GetPermissions()
		c.Set("permissions", user.Permissions)
		c.Next()
	}
}
