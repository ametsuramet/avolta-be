package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category represents a category in the blog application (optional)
type Category struct {
	Base
	Name        string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
}

func (u *Category) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.CategoryReponse{})
}
