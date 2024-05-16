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

func (u *Role) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.RoleReponse{})
}

func (m *Role) GetPermissions() {
	if m.IsSuperAdmin {
		database.DB.Find(&m.Permissions)
	}
}
