package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bank struct {
	Base
	Name string
	Code string
}

func (u *Bank) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Bank) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.BankResponse{
		ID:   m.ID,
		Name: m.Name,
		Code: m.Code,
	})
}
