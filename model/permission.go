package model

import (
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
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}
