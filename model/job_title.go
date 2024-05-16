package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobTitle struct {
	Base
	Name string
}

func (u *JobTitle) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m JobTitle) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.JobTitleReponse{
		ID:   m.ID,
		Name: m.Name,
	})
}
