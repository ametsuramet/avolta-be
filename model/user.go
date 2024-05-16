package model

import (
	"avolta/database"
	"avolta/object/resp"
	"avolta/util"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Base
	Email       string
	Password    string
	FullName    string
	Avatar      string
	IsAdmin     bool
	RoleID      *string
	Role        Role     `gorm:"foreignKey:RoleID"`
	Permissions []string `gorm:"-"`
}

func (u *User) GetPermissions() {
	u.Permissions = []string{}
	if u.RoleID != nil {
		role := Role{}
		database.DB.Find(&role, u.RoleID)
		role.GetPermissions()
		u.Role = role

		for _, v := range role.Permissions {
			u.Permissions = append(u.Permissions, v.Key)
		}
	}
}
func (u *User) CheckUserByEmail(email string) bool {
	count := int64(0)
	database.DB.Find(&u, "email = ?", email).Count(&count)
	return count > 0
}
func (u *User) GetUserByID() {
	count := int64(0)
	database.DB.Find(&u, "id = ?", u.ID).Count(&count)
}
func GetUserFromCtx(c *gin.Context) *User {
	user, ok := c.Get("user")
	if ok {
		authUser := user.(User)
		return &authUser
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	newPass, _ := util.HashPassword(u.Password)
	tx.Statement.SetColumn("password", newPass)

	return
}

func (m User) MarshalJSON() ([]byte, error) {
	m.GetPermissions()
	return json.Marshal(resp.UserResponse{
		FullName:    m.FullName,
		Permissions: m.Permissions,
		Email:       m.Email,
		Picture:     m.Avatar,
		RoleName:    m.Role.Name,
	})
}
