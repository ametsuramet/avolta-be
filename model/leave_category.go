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
	Absent      bool `gorm:"default:false; NOT NULL"`
	Sick        bool `gorm:"default:false; NOT NULL"`
}

func (u *LeaveCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m LeaveCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.LeaveCategoryResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Absent:      m.Absent,
		Sick:        m.Sick,
	})
}
