package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Incentive struct {
	Base
	IncentiveReportID   string                 `json:"incentive_report_id"`
	IncentiveReport     IncentiveReport        `gorm:"foreignKey:IncentiveReportID"`
	EmployeeID          string                 `json:"employee_id"`
	Employee            Employee               `gorm:"foreignKey:EmployeeID"`
	TotalSales          float64                `json:"total_sales"`
	TotalIncludedSales  float64                `json:"total_included_sales"`
	TotalExcludedSales  float64                `json:"total_excluded_sales"`
	TotalIncentive      float64                `json:"total_incentive"`
	TotalIncentiveBruto float64                `json:"total_incentive_bruto"`
	SickLeave           float64                `json:"sick_leave"`
	OtherLeave          float64                `json:"other_leave"`
	Absent              float64                `json:"absent"`
	Sales               []Sale                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	IncentiveShops      []IncentiveShop        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Status              string                 `json:"status" gorm:"type:enum('REQUEST','APPROVED', 'REJECTED', 'PAID', 'FINISHED');default:'REQUEST'"`
	Summaries           []ProductCategorySales `json:"summaries" gorm:"serializer:json;type:JSON"`
	CompanyID           string                 `json:"company_id" gorm:"not null"`
	Company             Company                `gorm:"foreignKey:CompanyID"`
}

func (u *Incentive) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Incentive) MapToResp() resp.IncentiveResponse {
	return resp.IncentiveResponse{
		ID:                  m.ID,
		IncentiveReportID:   m.IncentiveReportID,
		EmployeeID:          m.EmployeeID,
		EmployeeName:        m.Employee.FullName,
		TotalSales:          m.TotalSales,
		TotalIncludedSales:  m.TotalIncludedSales,
		TotalExcludedSales:  m.TotalExcludedSales,
		TotalIncentive:      m.TotalIncentive,
		TotalIncentiveBruto: m.TotalIncentiveBruto,
		SickLeave:           m.SickLeave,
		OtherLeave:          m.OtherLeave,
		Absent:              m.Absent,
		Summaries:           m.Summaries,
		Sales: func(items []Sale) []resp.SaleResponse {
			mapped := []resp.SaleResponse{}
			for _, v := range items {
				mapped = append(mapped, v.MapToResp())
			}
			return mapped
		}(m.Sales),
		IncentiveShops: func(items []IncentiveShop) []resp.IncentiveShopResponse {
			mapped := []resp.IncentiveShopResponse{}
			for _, v := range items {
				mapped = append(mapped, v.MapToResp())
			}
			return mapped
		}(m.IncentiveShops),
	}
}

func (m Incentive) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MapToResp())
}
