package handler

import (
	"avolta/model"
	"avolta/object/auth"
	"avolta/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Login(c *gin.Context) {

	input := auth.LoginRequest{}

	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error Binding"})
		return
	}

	var user model.User

	if exists := user.CheckUserByEmail(input.Email); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Exists"})
		return
	}
	fmt.Println(input.Password)
	fmt.Println(user.Password)
	if ok := util.CheckPasswordHash(input.Password, user.Password); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
		return
	}

	token, err := util.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
