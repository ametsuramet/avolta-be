package handler

import (
	"avolta/cmd"
	"avolta/model"
	"avolta/object/auth"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Login(c *gin.Context) {

	input := auth.LoginRequest{}

	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error Binding", "message": "Error Binding"})
		return
	}

	var user model.User

	if exists := user.CountSuperAdmin(); !exists {

		if err := user.CreateSuperAdmin(input.Email, input.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": "Error Generate Super Admin"})
			return
		}
		cmd.GenPermissions()
		cmd.GenBanks()
		cmd.GenAccounts()
		cmd.GenLeaveCategories()
		cmd.GenProductCategories()

	}

	if exists := user.CheckAdminByEmail(input.Email); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Exists", "message": "User Not Exists or not admin"})
		return
	}
	// fmt.Println(input.Password)
	// fmt.Println(user.Password)
	if ok := util.CheckPasswordHash(input.Password, user.Password); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password", "message": "Wrong Password"})
		return
	}

	token, err := util.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
func LoginUser(c *gin.Context) {

	input := auth.LoginRequest{}

	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error Binding", "message": "Error Binding"})
		return
	}

	var user model.User

	if exists := user.CheckUserByEmail(input.Email); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Not Exists", "message": "User Not Exists or not admin"})
		return
	}
	// fmt.Println(input.Password)
	// fmt.Println(user.Password)
	if ok := util.CheckPasswordHash(input.Password, user.Password); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password", "message": "Wrong Password"})
		return
	}

	token, err := util.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
