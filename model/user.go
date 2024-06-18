package model

import (
	"avolta/config"
	"avolta/database"
	"avolta/object/resp"
	svc "avolta/service"
	"avolta/util"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Base
	Email       string    `json:"email" gorm:"size:191;uniqueIndex"`
	Password    string    `json:"password"`
	FullName    string    `json:"full_name"`
	Avatar      string    `json:"avatar"`
	IsAdmin     bool      `json:"is_admin"`
	Token       *string   `json:"token"`
	Permissions []string  `gorm:"-"`
	EmployeeID  string    `json:"employee_id" gorm:"-"`
	Employee    Employee  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Companies   []Company `gorm:"many2many:user_companies;"`
	Roles       []Role    `gorm:"many2many:user_roles;"`
	VerifiedAt  *time.Time
	// RoleID      *string   `json:"role_id"`
}

func (u *User) GetPermissions(companyID string) {
	u.Permissions = []string{}
	// if u.RoleID != nil {
	// 	role := Role{}
	// 	database.DB.Preload("Permissions").Find(&role, u.RoleID)
	// 	role.GetPermissions()
	// 	u.Role = role

	// 	for _, v := range role.Permissions {
	// 		u.Permissions = append(u.Permissions, v.Key)
	// 	}
	// }
}
func (u *User) CheckAdminByEmail(email string) bool {
	count := int64(0)
	database.DB.Preload("Companies").Find(&u, "email = ? and is_admin = 1 and verified_at is not null", email).Count(&count)
	return count > 0
}

func (u *User) CountSuperAdmin() bool {
	count := int64(0)
	database.DB.Find(&u, "is_admin = 1").Count(&count)
	return count > 0
}
func (u *User) CreateSuperAdmin(email string, password string, fullName string) (string, error) {
	token := util.RandomString(20)
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		user := &User{
			Email:    email,
			FullName: fullName,
			Password: password,
			IsAdmin:  true,
			Token:    &token,
		}
		// CREATE FIRST SUPER ADMIN
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// // GENERATE STATIC DATA
		// role := Role{
		// 	Name:         "SUPERADMIN",
		// 	Description:  "Yes i'am superman",
		// 	IsSuperAdmin: true,
		// }

		// // CREATE SUPER ADMIN
		// if err := tx.Create(&role).Error; err != nil {
		// 	return err
		// }

		// if err := tx.Model(&user).Update("role_id", role.ID).Error; err != nil {
		// 	return err
		// }

		link := fmt.Sprintf("%s/verification/%s", config.App.Server.CmsURL, token)

		svc.MAIL.SetAddress(fullName, email)
		svc.MAIL.SetTemplate("template/layout.html", "template/new_user.html")
		err := svc.MAIL.SendEmail("Pendaftaran User", gin.H{
			"Name": fullName,
			"Link": link,
		}, []string{})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *User) CheckUserByEmail(email string) bool {
	count := int64(0)
	database.DB.Find(&u, "email = ?", email).Count(&count)
	return count > 0
}
func (u *User) GetUserByID(companyID string) {
	count := int64(0)
	database.DB.Preload("Roles", "company_id = ?", companyID).Find(&u, "id = ?", u.ID).Count(&count)
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
	// m.GetPermissions()
	return json.Marshal(resp.UserResponse{
		ID:          m.ID,
		FullName:    m.FullName,
		Permissions: m.Permissions,
		Email:       m.Email,
		Picture:     m.Avatar,
		// RoleName:     m.Role.Name,
		// RoleID:       m.Role.ID,
		IsAdmin:      m.IsAdmin,
		EmployeeID:   m.Employee.ID,
		EmployeeName: m.Employee.FullName,
	})
}
