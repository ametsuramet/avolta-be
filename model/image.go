package model

import (
	"encoding/json"

	"avolta/object/resp"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Image represents an image associated with a blog post (optional)
type Image struct {
	Base
	Filename string `gorm:"not null"`
	Path     string `gorm:"not null"`
	Caption  string
}

func (u *Image) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Image) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.ImageResponse{})
}
