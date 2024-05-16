package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenPermission struct {
	Base
	Name string
}

func (u *GenPermission) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m GenPermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.GenPermissionReponse{})
}
