package model

import (
	"avolta/database"
	"avolta/object/resp"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayRollReport struct {
	Base
	Description           string              `json:"description"`
	ReportNumber          string              `json:"report_number"`
	UserID                string              `json:"user_id"`
	User                  User                `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StartDate             time.Time           `json:"start_date"`
	EndDate               time.Time           `json:"end_date"`
	Status                string              `json:"status" gorm:"type:enum('DRAFT',  'PROCESSING', 'FINISHED', 'CANCELED');default:'DRAFT'"`
	Items                 []PayRollReportItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GrandTotalTakeHomePay float64             `json:"grand_total_take_home_pay"`
	CompanyID             string              `json:"company_id"`
	Company               Company             `gorm:"foreignKey:CompanyID"`
}

func (u *PayRollReport) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m *PayRollReport) MapToResp() resp.PayRollReportResponse {
	return resp.PayRollReportResponse{
		ID:                    m.ID,
		Description:           m.Description,
		ReportNumber:          m.ReportNumber,
		UserID:                m.UserID,
		UserName:              m.User.FullName,
		StartDate:             m.StartDate,
		EndDate:               m.EndDate,
		Status:                m.Status,
		GrandTotalTakeHomePay: m.GrandTotalTakeHomePay,
		Items: func(items []PayRollReportItem) []resp.PayRollReportItemResponse {
			mapped := []resp.PayRollReportItemResponse{}
			for _, v := range items {
				mapped = append(mapped, v.MapToResp())
			}
			return mapped
		}(m.Items),
	}
}
func (m PayRollReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MapToResp())
}

func (m PayRollReport) AddItem(payrollID string) error {
	payroll := PayRoll{}
	if err := database.DB.Preload("Employee").Preload("Employee.Bank").Find(&payroll, "id = ?", payrollID).Error; err != nil {
		return err
	}
	if payroll.PayRollReportItemID != nil {
		err := errors.New("payroll has reported")
		return err
	}
	item := &PayRollReportItem{
		TotalTakeHomePay:       payroll.TakeHomePay,
		TotalReimbursement:     payroll.TotalReimbursement,
		EmployeeName:           payroll.Employee.FullName,
		EmployeeID:             payroll.Employee.ID,
		EmployeeEmail:          payroll.Employee.Email,
		EmployeePhone:          payroll.Employee.Phone,
		EmployeeIdentityNumber: payroll.Employee.EmployeeIdentityNumber,
		BankAccountNumber:      payroll.Employee.BankAccountNumber,
		BankName:               payroll.Employee.Bank.Name,
		BankCode:               payroll.Employee.Bank.Code,
		BankID:                 payroll.Employee.Bank.ID,
		PayRollReportID:        m.ID,
	}
	if err := database.DB.Create(&item).Error; err != nil {
		return err
	}

	if err := database.DB.Model(&payroll).Update("pay_roll_report_item_id", item.ID).Error; err != nil {
		return err
	}

	return m.UpdateGrandTotal()
}

func (m PayRollReport) DeleteItem(payrollID string) error {
	payroll := PayRoll{}
	if err := database.DB.Find(&payroll, "id = ? and pay_roll_report_item_id is not null", payrollID).Error; err != nil {
		return err
	}
	payrollReportItem := PayRollReportItem{}
	if err := database.DB.Find(&payrollReportItem, "pay_roll_report_item_id = ?", payroll.PayRollReportItemID).Error; err != nil {
		return err
	}
	if err := database.DB.Model(&payroll).Update("pay_roll_report_item_id", nil).Error; err != nil {
		return err
	}
	if err := database.DB.Model(&payrollReportItem).Delete(&payrollReportItem, "id = ?", payrollReportItem.ID).Error; err != nil {
		return err
	}

	return m.UpdateGrandTotal()
}

func (m *PayRollReport) GetItems() {
	items := []PayRollReportItem{}
	database.DB.Find(&items, "pay_roll_report_id = ?", m.ID)
	m.Items = items
}
func (m PayRollReport) UpdateGrandTotal() error {
	m.GetItems()
	total := float64(0)
	for _, v := range m.Items {
		total += v.TotalTakeHomePay + v.TotalReimbursement
	}

	m.GrandTotalTakeHomePay = total
	if err := database.DB.Model(&m).Update("grand_total_take_home_pay", total).Error; err != nil {
		return err
	}
	return nil
}
