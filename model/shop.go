package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shop struct {
	Base
	Name string
}

func (u *Shop) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Shop) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.ShopReponse{})
}
