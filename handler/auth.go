package handler

import (
	"avolta/cmd"
	"avolta/database"
	"avolta/model"
	"avolta/object/auth"
	"avolta/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func CreateCompany(c *gin.Context) {
	var data model.Company

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	getUser, _ := c.Get("user")
	user := getUser.(model.User)

	data.UserID = user.ID

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Create(&data).Error; err != nil {
			return err
		}
		// GENERATE STATIC DATA
		role := model.Role{
			Name:         "SUPERADMIN",
			Description:  "Yes i'am superman",
			IsSuperAdmin: true,
			CompanyID:    data.ID,
		}

		// CREATE SUPER ADMIN
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		tx.Model(&user).Association("Companies").Append(&data)
		tx.Model(&user).Association("Roles").Append(&role)
		cmd.GenAccounts(tx, []string{"", data.ID})
		return nil
	})

	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data Company Created", gin.H{"last_id": data.ID}, nil)
}
func Verification(c *gin.Context) {
	token := c.Params.ByName("token")
	var user model.User
	if err := database.DB.Find(&user, "token = ?", token).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if user.VerifiedAt != nil {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User, %s, telah Terverfifikasi", user.FullName)})
		return
	}
	now := time.Now()
	user.VerifiedAt = &now

	if err := database.DB.Updates(&user).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Selamat Bergabung, %s, Terimakasih atas kepercayaan anda", user.FullName)})
}
func Register(c *gin.Context) {
	input := auth.RegisterRequest{}

	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error Binding", "message": "Error Binding"})
		return
	}

	var user model.User
	// if exists := user.CountSuperAdmin(); !exists {
	_, err := user.CreateSuperAdmin(input.Email, input.Password, input.FullName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	// 	cmd.GenPermissions()
	// 	cmd.GenBanks()
	// 	cmd.GenAccounts()
	// 	cmd.GenLeaveCategories()
	// 	cmd.GenProductCategories()

	// }
	c.JSON(http.StatusOK, gin.H{"message": "Silahkan cek email anda"})
}

func Login(c *gin.Context) {

	input := auth.LoginRequest{}

	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error Binding", "message": "Error Binding"})
		return
	}

	var user model.User

	// if exists := user.CountSuperAdmin(); !exists {

	// 	if err := user.CreateSuperAdmin(input.Email, input.Password); err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "message": "Error Generate Super Admin"})
	// 		return
	// 	}
	// 	cmd.GenPermissions()
	// 	cmd.GenBanks()
	// 	cmd.GenAccounts()
	// 	cmd.GenLeaveCategories()
	// 	cmd.GenProductCategories()

	// }

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

	c.JSON(http.StatusOK, gin.H{"token": token, "companies": user.Companies})

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
