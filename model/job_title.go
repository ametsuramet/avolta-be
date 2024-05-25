package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobTitle struct {
	Base
	Name        string
	Description string
}

func (u *JobTitle) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m JobTitle) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.JobTitleResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
	})
}
