package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name string
}

func (u *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.ProductReponse{})
}
