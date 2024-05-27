package model

import (
	"avolta/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	Base
	Name     string    `json:"name"`
	Products []Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

type ProductCategorySales struct {
	ID               string  `json:"id"`
	ShopID           string  `json:"shop_id"`
	Name             string  `json:"name"`
	ShopName         string  `json:"shop_name"`
	Total            float64 `json:"total"`
	ComissionPercent float64 `json:"commission_percent"`
	TotalComission   float64 `json:"total_comission"`
}

func (m *ProductCategorySales) GetIncentive() {

	prodCat := ProductCategory{}
	incentiveSetting := IncentiveSetting{}
	shop := Shop{}
	database.DB.Find(&incentiveSetting, "shop_id = ? and product_category_id = ?", m.ShopID, m.ID)
	database.DB.Select("name").Find(&prodCat, "id = ?", m.ID)
	database.DB.Select("name").Find(&shop, "id = ?", m.ShopID)

	// if sick > incentiveSetting.SickLeaveThreshold {
	// 	return 0
	// }
	// if leave > incentiveSetting.OtherLeaveThreshold {
	// 	return 0
	// }
	// if absent > incentiveSetting.AbsentThreshold {
	// 	return 0
	// }
	commissionPercent := float64(0)
	if m.Total > incentiveSetting.MinimumSalesTarget {
		commissionPercent = incentiveSetting.MinimumSalesCommission
	}
	if m.Total > incentiveSetting.MaximumSalesTarget {
		commissionPercent = incentiveSetting.MaximumSalesCommission
	}
	m.ComissionPercent = commissionPercent
	m.Name = prodCat.Name
	m.ShopName = shop.Name
	m.TotalComission = m.Total * commissionPercent
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
