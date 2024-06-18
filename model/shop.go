package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shop struct {
	Base
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Code           string          `json:"code"`
	Sales          []Sale          `json:"sales" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IncentiveShops []IncentiveShop `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CompanyID      string          `json:"company_id" gorm:"not null"`
	Company        Company         `gorm:"foreignKey:CompanyID"`
}

func (u *Shop) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Shop) MapToResp() resp.ShopResponse {
	return resp.ShopResponse{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
	}
}
func (m Shop) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MapToResp())
}
