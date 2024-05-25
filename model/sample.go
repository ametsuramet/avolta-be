package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sample struct {
	Base
	Name string
}

func (u *Sample) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Sample) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.SampleResponse{})
}
