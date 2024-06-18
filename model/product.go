package model

import (
	"avolta/object/resp"
	"avolta/util"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name              string          `json:"name"`
	SKU               string          `json:"sku"`
	Barcode           string          `json:"barcode"`
	SellingPrice      float64         `json:"selling_price"`
	ProductCategoryID *string         `json:"product_category_id"`
	ProductCategory   ProductCategory `json:"product_category" gorm:"foreignKey:ProductCategoryID"`
	Sale              []Sale          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CompanyID         string          `json:"company_id" gorm:"not null"`
	Company           Company         `gorm:"foreignKey:CompanyID"`
}

func (u *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.ProductResponse{
		ID:                  m.ID,
		Name:                m.Name,
		SKU:                 m.SKU,
		Barcode:             m.Barcode,
		SellingPrice:        m.SellingPrice,
		ProductCategoryID:   util.SavedString(m.ProductCategoryID),
		ProductCategoryName: m.ProductCategory.Name,
	})
}
