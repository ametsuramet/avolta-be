package model

import (
	"avolta/object/resp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Incentive struct {
	Base
	StartDate          time.Time       `json:"start_date"`
	EndDate            time.Time       `json:"end_date"`
	EmployeeID         string          `json:"employee_id"`
	Employee           Employee        `gorm:"foreignKey:EmployeeID"`
	TotalSales         float64         `json:"total_sales"`
	TotalIncludedSales float64         `json:"total_included_sales"`
	TotalExcludedSales float64         `json:"total_excluded_sales"`
	TotalIncentive     float64         `json:"total_incentive"`
	SickLeave          float64         `json:"sick_leave"`
	OtherLeave         float64         `json:"other_leave"`
	Absent             float64         `json:"absent"`
	IncentiveShops     []IncentiveShop `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (u *Incentive) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Incentive) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.IncentiveResponse{})
}
