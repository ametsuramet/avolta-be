package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	Base
	Name     string    `json:"name"`
	Products []Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (u *ProductCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

// func (m ProductCategory) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(resp.ProductCategoryResponse{})
// }
