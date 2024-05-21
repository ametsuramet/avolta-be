package middleware

import (
	"avolta/config"
	"avolta/model"
	"avolta/object/auth"
	"avolta/util"
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

		token1, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.App.Server.SecretKey), nil
		})

		fmt.Println(token1.Claims)

		token, err := jwt.ParseWithClaims(reqToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.App.Server.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			fmt.Println("Authorized user:", claims.Id)
			var user = model.User{Base: model.Base{ID: claims.Id}}
			user.GetUserByID()
			// fmt.Println("user", user)
			timezone := c.GetHeader("timezone")
			if timezone == "" {
				timezone = "Asia/Jakarta"
			}
			c.Set("user", user)
			c.Set("timezone", timezone)
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
		// fmt.Println("PERMISSIONS", user.Permissions)
		c.Set("permissions", user.Permissions)
		c.Next()
	}
}
func PermissionMiddleware(userPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		getPermissions, _ := c.Get("permissions")
		permissions := getPermissions.([]string)
		getUser, _ := c.Get("user")
		user := getUser.(model.User)
		user.GetPermissions()
		if user.Role.IsSuperAdmin {
			c.Next()
			return
		}

		if util.Contains(permissions, userPermission) {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no " + userPermission + " permission exists"})
		c.Abort()
	}
}
