package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncentiveSetting struct {
	Base
	ShopID                 string          `json:"shop_id"`
	Shop                   Shop            `json:"shop" gorm:"foreignKey:ShopID"`
	ProductCategoryID      string          `json:"product_category_id"`
	ProductCategory        ProductCategory `json:"product_category" gorm:"foreignKey:ProductCategoryID"`
	MinimumSalesTarget     float64         `json:"minimum_sales_target"`
	MaximumSalesTarget     float64         `json:"maximum_sales_target"`
	MinimumSalesCommission float64         `json:"minimum_sales_commission"`
	MaximumSalesCommission float64         `json:"maximum_sales_commission"`
	SickLeaveThreshold     float64         `json:"sick_leave_threshold"`
	OtherLeaveThreshold    float64         `json:"other_leave_threshold"`
	AbsentThreshold        float64         `json:"absent_threshold"`
}

func (u *IncentiveSetting) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m IncentiveSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.IncentiveSettingResponse{})
}
