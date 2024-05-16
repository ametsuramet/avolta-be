package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Leave struct {
	Base
	Name string
}

func (u *Leave) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Leave) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.LeaveReponse{})
}
