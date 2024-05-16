package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	Base
	Parent      *Organization
	ParentId    string     `gorm:"TYPE:integer REFERENCES organizations"`
	Name        string     `gorm:"size:30"`
	Description string     `gorm:"size:100"`
	Employee    []Employee `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Organization) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.OrganizationReponse{})
}
