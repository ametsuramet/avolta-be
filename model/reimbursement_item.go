package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReimbursementItem struct {
	Base
	Amount          float64       `json:"amount" form:"amount"`
	Notes           string        `json:"notes" form:"notes"`
	ReimbursementID string        `json:"reimbursement_id"`
	Reimbursement   Reimbursement `gorm:"foreignKey:ReimbursementID" json:"-"`
	Files           string        `json:"files" gorm:"default:'[]'"`
	Attachments     []string      `json:"attachments" gorm:"-"`
}

func (u *ReimbursementItem) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m ReimbursementItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.ReimbursementItemReponse{})
}
