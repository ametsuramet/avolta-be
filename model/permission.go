package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	Base
	Name          string `json:"name"`
	Key           string `json:"key" gorm:"unique"`
	Description   string `json:"description"`
	DescriptionEn string `json:"description_en"`
	Group         string `json:"group"`
	IsDefault     bool   `json:"is_default"`
	IsActive      bool   `json:"is_active"`
	IsSpecial     bool   `json:"is_special"`
}

func (u *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.PermissionReponse{})
}
