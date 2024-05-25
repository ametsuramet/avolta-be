package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SaleReceipt struct {
	Base
	Name  string
	Sales []Sale  `json:"sales" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Total float64 `json:"total"`
}

func (u *SaleReceipt) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m SaleReceipt) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.SaleReceiptResponse{})
}
