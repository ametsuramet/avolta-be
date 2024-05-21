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
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	FullName    string   `json:"full_name"`
	Avatar      string   `json:"avatar"`
	IsAdmin     bool     `json:"is_admin"`
	RoleID      *string  `json:"role_id"`
	Role        Role     `json:"role" gorm:"foreignKey:RoleID"`
	Permissions []string `gorm:"-"`
	EmployeeID  string   `json:"employee_id" gorm:"-"`
	Employee    Employee `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *User) GetPermissions() {
	u.Permissions = []string{}
	if u.RoleID != nil {
		role := Role{}
		database.DB.Preload("Permissions").Find(&role, u.RoleID)
		role.GetPermissions()
		u.Role = role

		for _, v := range role.Permissions {
			u.Permissions = append(u.Permissions, v.Key)
		}
	}
}
func (u *User) CheckAdminByEmail(email string) bool {
	count := int64(0)
	database.DB.Find(&u, "email = ? and is_admin = 1", email).Count(&count)
	return count > 0
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
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	newPass, _ := util.HashPassword(u.Password)
	tx.Statement.SetColumn("password", newPass)

	return
}

func (m User) MarshalJSON() ([]byte, error) {
	m.GetPermissions()
	return json.Marshal(resp.UserResponse{
		ID:           m.ID,
		FullName:     m.FullName,
		Permissions:  m.Permissions,
		Email:        m.Email,
		Picture:      m.Avatar,
		RoleName:     m.Role.Name,
		RoleID:       m.Role.ID,
		IsAdmin:      m.IsAdmin,
		EmployeeID:   m.Employee.ID,
		EmployeeName: m.Employee.FullName,
	})
}
