package model

import (
	"avolta/database"
	"avolta/object/resp"
	"avolta/util"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sale struct {
	Base
	Date                time.Time   `json:"date"`
	Code                string      `json:"code"`
	SaleReceiptID       *string     `json:"sale_receipt_id"`
	SaleReceipt         SaleReceipt `json:"sale_receipt" gorm:"foreignKey:SaleReceiptID"`
	ProductID           string      `json:"product_id"`
	Product             Product     `json:"product" gorm:"foreignKey:ProductID"`
	ShopID              string      `json:"shop_id"`
	Shop                Shop        `json:"shop" gorm:"foreignKey:ShopID"`
	Qty                 float64     `json:"qty"`
	Price               float64     `json:"price"`
	SubTotal            float64     `json:"sub_total"`
	Discount            float64     `json:"discount"`
	DiscountAmount      float64     `json:"discount_amount"`
	Total               float64     `json:"total"`
	EmployeeID          string      `json:"employee_id"`
	Employee            Employee    `gorm:"foreignKey:EmployeeID"`
	IncentiveID         *string     `json:"incentive_id"`
	Incentive           Incentive   `gorm:"foreignKey:IncentiveID"`
	IsIncentiveExcluded bool        `json:"is_incentive_excluded"`
	IsIncentiveIncluded bool        `json:"is_incentive_included"`
}

func (u *Sale) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Sale) GetIncentive(shopId string, sick float64, leave float64, absent float64) float64 {
	product := Product{}
	database.DB.Select("product_category_id").Find(&product, "id = ?", m.ProductID)

	incentiveSetting := IncentiveSetting{}
	database.DB.Find(&incentiveSetting, "shop_id = ? and product_category_id = ?", shopId, product.ProductCategoryID)

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
	fmt.Println(incentiveSetting.MinimumSalesTarget, incentiveSetting.MaximumSalesCommission, commissionPercent)
	// ADD OTHER CRITERIA HERE
	return m.Total * commissionPercent
}

func (m Sale) MapToResp() resp.SaleResponse {
	employeePicture := ""
	if m.Employee.Picture.Valid {
		employeePicture = util.GetURL(m.Employee.Picture.String)
	}
	return resp.SaleResponse{
		ID:              m.ID,
		Date:            m.Date.Format("2006-01-02 15:04:05"),
		Code:            m.Code,
		ProductID:       m.ProductID,
		ProductName:     m.Product.Name,
		ProductSKU:      m.Product.SKU,
		ShopID:          m.ShopID,
		ShopName:        m.Shop.Name,
		Qty:             m.Qty,
		Price:           m.Price,
		SubTotal:        m.SubTotal,
		Discount:        m.Discount,
		DiscountAmount:  m.DiscountAmount,
		Total:           m.Total,
		EmployeeID:      m.EmployeeID,
		EmployeeName:    m.Employee.FullName,
		EmployeePicture: employeePicture,
		IncentiveID:     m.IncentiveID,
	}
}
func (m Sale) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MapToResp())
}
