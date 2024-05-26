package model

import (
	"avolta/database"
	"avolta/object/resp"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncentiveReport struct {
	Base
	Description  string      `json:"description"`
	ReportNumber string      `json:"report_number"`
	UserID       string      `json:"user_id"`
	User         User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StartDate    time.Time   `json:"start_date"`
	EndDate      time.Time   `json:"end_date"`
	Shops        []Shop      `gorm:"many2many:incentive_report_shops;"`
	Status       string      `json:"status" gorm:"type:enum('DRAFT',  'PROCESSING', 'FINISHED', 'CANCELED');default:'DRAFT'"`
	Incentives   []Incentive `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type IncentiveReportReq struct {
	Description  string    `json:"description"`
	ReportNumber string    `json:"report_number"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	ShopIDs      []string  `json:"shop_ids"`
}

func (u *IncentiveReport) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m IncentiveReport) MapToResp() resp.IncentiveReportResponse {
	return resp.IncentiveReportResponse{
		ID:           m.ID,
		Description:  m.Description,
		ReportNumber: m.ReportNumber,
		UserID:       m.UserID,
		UserName:     m.User.FullName,
		StartDate:    m.StartDate,
		EndDate:      m.EndDate,
		Status:       m.Status,
		Incentives: func(items []Incentive) []resp.IncentiveResponse {
			mapped := []resp.IncentiveResponse{}
			for _, v := range items {
				mapped = append(mapped, v.MapToResp())
			}
			return mapped
		}(m.Incentives),
		Shops: func(items []Shop) []resp.ShopResponse {
			mapped := []resp.ShopResponse{}
			for _, v := range items {
				mapped = append(mapped, v.MapToResp())
			}
			return mapped
		}(m.Shops),
	}
}
func (m IncentiveReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MapToResp())
}

func (m *IncentiveReport) AddEmployee(c *gin.Context) error {
	input := struct {
		EmployeeID string `json:"employee_id" binding:"required"`
	}{}

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		return err
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		employee := Employee{}
		if err := tx.Find(&employee, "id = ?", input.EmployeeID).Error; err != nil {
			return err
		}

		incentive := Incentive{
			IncentiveReportID: m.ID,
			EmployeeID:        input.EmployeeID,
			SickLeave:         0,
			OtherLeave:        0,
			Absent:            0,
		}

		if err := tx.Create(&incentive).Error; err != nil {
			return err
		}

		grandTotalSales := float64(0)
		grandTotalComission := float64(0)
		for _, v := range m.Shops {

			productCatSales := []ProductCategorySales{}

			tx.Raw(`SELECT  sum(sales.total) as total, p.product_category_id id, pc.name , sales.shop_id  FROM sales
			JOIN products p ON p.id = sales.product_id 
			JOIN product_categories pc ON pc.id = p.product_category_id 
			WHERE (date BETWEEN ? and ?) AND (employee_id = ? and shop_id = ?) AND sales.deleted_at IS NULL
			GROUP by p.product_category_id, pc.name, sales.shop_id ;`, m.StartDate, m.EndDate, input.EmployeeID, v.ID).Scan(&productCatSales)

			// util.LogJson(productCatSales)
			totalSales := float64(0)
			totalComission := float64(0)
			for _, v := range productCatSales {
				totalSales += v.Total
				commision := v.GetIncentive()
				totalComission += commision
			}

			// CREATE INCENTIVE SHOP DATA
			incentiveShop := IncentiveShop{
				ShopID:         v.ID,
				IncentiveID:    incentive.ID,
				TotalSales:     totalSales,
				TotalIncentive: totalComission,
			}

			if err := tx.Create(&incentiveShop).Error; err != nil {
				return err
			}

			grandTotalSales += totalSales
			grandTotalComission += totalComission
		}
		incentive.TotalSales = grandTotalSales
		incentive.TotalIncentive = grandTotalComission
		if err := tx.Model(&incentive).Updates(&incentive).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
