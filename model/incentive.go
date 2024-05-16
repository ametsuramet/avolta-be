package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Incentive struct {
	Base
	Name string
}

func (u *Incentive) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Incentive) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.IncentiveReponse{})
}
