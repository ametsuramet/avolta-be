package model

import (
	"avolta/database"
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	Base
	Name         string
	Description  string
	IsSuperAdmin bool
	Permissions  []Permission `gorm:"many2many:role_permissions;"`
	User         []User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type RoleReq struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	IsSuperAdmin bool     `json:"is_super_admin"`
	Permissions  []string `json:"permissions"`
}

func (u *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Role) MarshalJSON() ([]byte, error) {
	permissions := []string{}
	for _, v := range m.Permissions {
		permissions = append(permissions, v.Key)
	}
	return json.Marshal(resp.RoleResponse{
		ID:           m.ID,
		Name:         m.Name,
		Description:  m.Description,
		IsSuperAdmin: m.IsSuperAdmin,
		Permissions:  permissions,
	})
}

func (m *Role) GetPermissions() {
	if m.IsSuperAdmin {
		database.DB.Find(&m.Permissions)
	}

}
