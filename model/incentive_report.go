package model

import (
	"avolta/database"
	"avolta/object/resp"
	"encoding/json"
	"errors"
	"fmt"
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
	Shops        []Shop      `gorm:"many2many:incentive_report_shops;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status       string      `json:"status" gorm:"type:enum('DRAFT',  'PROCESSING', 'FINISHED', 'CANCELED');default:'DRAFT'"`
	Incentives   []Incentive `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CompanyID    string      `json:"company_id" gorm:"not null"`
	Company      Company     `gorm:"foreignKey:CompanyID"`
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

func (m *IncentiveReport) UpdateIncentive(c *gin.Context) error {
	input := struct {
		SickLeave  float64 `json:"sick_leave" `
		OtherLeave float64 `json:"other_leave" `
		Absent     float64 `json:"absent" `
	}{}

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		return err
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		setting := Setting{}
		if err := tx.First(&setting).Error; err != nil {
			return err
		}

		incentiveId := c.Params.ByName("incentiveId")
		incentive := Incentive{}
		if err := tx.Preload("IncentiveShops").Find(&incentive, "id = ?", incentiveId).Error; err != nil {
			return err
		}

		incentive.SickLeave = input.SickLeave
		incentive.OtherLeave = input.OtherLeave
		incentive.Absent = input.Absent

		incentiveIds := []string{}
		for _, v := range incentive.IncentiveShops {
			incentiveIds = append(incentiveIds, v.ID)
		}

		grandTotalSales := float64(0)
		grandTotalComission := float64(0)
		grandTotalComissionBruto := float64(0)
		allProductCats := []ProductCategorySales{}
		for _, shop := range m.Shops {

			totalSales := float64(0)
			totalComission := float64(0)
			totalComissionBruto := float64(0)
			// GET SALES per shops
			sales := []Sale{}
			if err := tx.Preload("Product").Where("`date` BETWEEN ? and ?", m.StartDate, m.EndDate).Find(&sales, "employee_id = ? and shop_id = ? and incentive_id = ?", incentive.EmployeeID, shop.ID, incentive.ID).Error; err != nil {
				return err
			}

			if len(sales) == 0 {
				continue
			}

			sumCat := make(map[string]float64)
			for _, v := range sales {
				sumCat[*v.Product.ProductCategoryID] += v.Total

			}
			productCats := []ProductCategorySales{}
			for id, total := range sumCat {
				productCat := ProductCategorySales{}
				productCat.ID = id
				productCat.ShopID = shop.ID
				productCat.Total = total
				productCat.GetIncentive()
				totalSales += total
				totalComission += productCat.TotalComission
				totalComissionBruto += productCat.TotalComission
				if incentive.SickLeave > setting.IncentiveSickLeaveThreshold {
					totalComission = 0
				}
				if incentive.OtherLeave > setting.IncentiveOtherLeaveThreshold {
					totalComission = 0
				}
				if incentive.Absent > setting.IncentiveAbsentThreshold {
					totalComission = 0
				}
				fmt.Println("SETTING", incentive.SickLeave, "=>", setting.IncentiveSickLeaveThreshold)
				fmt.Println("TOTAL COMMISION", totalComission)
				productCats = append(productCats, productCat)
				if incentive.SickLeave == 0 {
					if err := tx.Model(&incentive).Update("sick_leave", 0).Error; err != nil {
						return err
					}
				}
				if incentive.OtherLeave == 0 {
					if err := tx.Model(&incentive).Update("other_leave", 0).Error; err != nil {
						return err
					}
				}
				if incentive.Absent == 0 {
					if err := tx.Model(&incentive).Update("absent", 0).Error; err != nil {
						return err
					}
				}
			}

			incentiveShop := IncentiveShop{}

			if err := tx.Where("id in (?)", incentiveIds).Find(&incentiveShop, "shop_id = ? and incentive_id = ?", shop.ID, incentive.ID).Error; err != nil {
				return err
			}
			// CREATE INCENTIVE SHOP DATA

			incentiveShop.TotalSales = totalSales
			incentiveShop.TotalIncentive = totalComission
			incentiveShop.TotalIncentiveBruto = totalComissionBruto
			incentiveShop.Summaries = productCats

			if err := tx.Model(&incentiveShop).Updates(&incentiveShop).Error; err != nil {
				return err
			}

			if totalComission == 0 {
				if err := tx.Model(&incentiveShop).Update("total_incentive", 0).Error; err != nil {
					return err
				}
			}
			if totalComissionBruto == 0 {
				if err := tx.Model(&incentiveShop).Update("total_comission_bruto", 0).Error; err != nil {
					return err
				}
			}

			grandTotalSales += totalSales
			grandTotalComission += totalComission
			grandTotalComissionBruto += totalComissionBruto
			allProductCats = append(allProductCats, productCats...)
			// UPDATE SALES

		}
		incentive.TotalSales = grandTotalSales
		incentive.TotalIncentive = grandTotalComission
		incentive.TotalIncentiveBruto = grandTotalComissionBruto
		incentive.Summaries = allProductCats
		if err := tx.Model(&incentive).Updates(&incentive).Error; err != nil {
			return err
		}

		if grandTotalComission == 0 {
			if err := tx.Model(&incentive).Update("total_incentive", 0).Error; err != nil {
				return err
			}
		}
		if grandTotalComissionBruto == 0 {
			if err := tx.Model(&incentive).Update("total_comission_bruto", 0).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}
func (m *IncentiveReport) AddEmployee(c *gin.Context) error {
	input := struct {
		EmployeeID string `json:"employee_id" binding:"required"`
	}{}

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		return err
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		setting := Setting{}
		if err := tx.First(&setting).Error; err != nil {
			return err
		}

		employee := Employee{}
		if err := tx.Find(&employee, "id = ?", input.EmployeeID).Error; err != nil {
			return err
		}

		count := int64(0)
		tx.Model(&Incentive{}).Where("employee_id = ? AND incentive_report_id =?", input.EmployeeID, m.ID).Count(&count)

		if count > 0 {
			err := errors.New("data insentive sudah ada")
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
		grandTotalComissionBruto := float64(0)
		allProductCats := []ProductCategorySales{}
		for _, shop := range m.Shops {

			totalSales := float64(0)
			totalComission := float64(0)
			totalComissionBruto := float64(0)
			// GET SALES per shops
			sales := []Sale{}
			if err := tx.Preload("Product").Where("`date` BETWEEN ? and ?", m.StartDate, m.EndDate).Find(&sales, "employee_id = ? and shop_id = ? and incentive_id is null", input.EmployeeID, shop.ID).Error; err != nil {
				return err
			}

			if len(sales) == 0 {
				continue
			}

			sumCat := make(map[string]float64)
			for _, v := range sales {
				sumCat[*v.Product.ProductCategoryID] += v.Total

			}
			productCats := []ProductCategorySales{}
			for id, total := range sumCat {
				productCat := ProductCategorySales{}
				productCat.ID = id
				productCat.ShopID = shop.ID
				productCat.Total = total
				productCat.GetIncentive()
				totalSales += total
				totalComission += productCat.TotalComission
				totalComissionBruto += productCat.TotalComission

				if incentive.SickLeave > setting.IncentiveSickLeaveThreshold {
					totalComission = 0
				}
				if incentive.OtherLeave > setting.IncentiveOtherLeaveThreshold {
					totalComission = 0
				}
				if incentive.Absent > setting.IncentiveAbsentThreshold {
					totalComission = 0
				}

				productCats = append(productCats, productCat)
			}

			// CREATE INCENTIVE SHOP DATA
			incentiveShop := IncentiveShop{
				ShopID:              shop.ID,
				IncentiveID:         incentive.ID,
				TotalSales:          totalSales,
				TotalIncentive:      totalComission,
				TotalIncentiveBruto: totalComissionBruto,
				Summaries:           productCats,
			}

			if err := tx.Create(&incentiveShop).Error; err != nil {
				return err
			}

			grandTotalSales += totalSales
			grandTotalComission += totalComission
			grandTotalComissionBruto += totalComissionBruto
			allProductCats = append(allProductCats, productCats...)
			// UPDATE SALES
			for _, v := range sales {
				v.IncentiveID = &incentive.ID
				tx.Model(&v).Updates(&v)
			}
		}
		incentive.TotalSales = grandTotalSales
		incentive.TotalIncentive = grandTotalComission
		incentive.TotalIncentiveBruto = grandTotalComissionBruto
		incentive.Summaries = allProductCats
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
