package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaveCategory struct {
	Base
	Name        string
	Description string
}

func (u *LeaveCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m LeaveCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.LeaveCategoryReponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
	})
}
