package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	Base
	Name string
}

func (u *ProductCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m ProductCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.ProductCategoryReponse{})
}
