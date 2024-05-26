package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncentiveShop struct {
	Base
	ShopID             string    `json:"shop_id"`
	Shop               Shop      `json:"shop" gorm:"foreignKey:ShopID"`
	IncentiveID        string    `json:"incentive_id"`
	Incentive          Incentive `gorm:"foreignKey:IncentiveID"`
	TotalSales         float64   `json:"total_sales"`
	TotalIncludedSales float64   `json:"total_included_sales"`
	TotalExcludedSales float64   `json:"total_excluded_sales"`
	TotalIncentive     float64   `json:"total_incentive"`
}

func (u *IncentiveShop) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m IncentiveShop) MapToResp() resp.IncentiveShopResponse {
	return resp.IncentiveShopResponse{
		ID:                 m.ID,
		ShopID:             m.ShopID,
		ShopName:           m.Shop.Name,
		IncentiveID:        m.IncentiveID,
		TotalSales:         m.TotalSales,
		TotalIncludedSales: m.TotalIncludedSales,
		TotalExcludedSales: m.TotalExcludedSales,
		TotalIncentive:     m.TotalIncentive,
	}
}
func (m IncentiveShop) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MapToResp())
}
