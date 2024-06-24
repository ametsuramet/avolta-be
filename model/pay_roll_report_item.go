package model

import (
	"avolta/database"
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayRollReportItem struct {
	Base
	TotalTakeHomePay       float64       `json:"total_take_home_pay"`
	TotalReimbursement     float64       `json:"total_reimbursement"`
	EmployeeName           string        `json:"employee_name"`
	EmployeePhone          string        `json:"employee_phone"`
	EmployeeEmail          string        `json:"employee_email"`
	EmployeeID             string        `json:"employee_id" gorm:"size:191"`
	Employee               Employee      `gorm:"foreignKey:EmployeeID;" json:"-"`
	EmployeeIdentityNumber string        `json:"employee_identity_number"`
	BankAccountNumber      string        `json:"bank_account_number"`
	BankName               string        `json:"bank_name"`
	BankCode               string        `json:"bank_code"`
	BankID                 string        `json:"bank_id"`
	Bank                   Bank          `gorm:"foreignKey:BankID" json:"-"`
	PayRollReportID        string        `json:"pay_roll_report_id"`
	PayRollReport          PayRollReport `gorm:"foreignKey:PayRollReportID" json:"-"`
	Status                 string        `json:"status" gorm:"type:enum('PROCESSING', 'FINISHED', 'PAID');default:'PROCESSING'"`
	PayRoll                PayRoll       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CompanyID              string        `json:"company_id" gorm:"not null"`
	Company                Company       `gorm:"foreignKey:CompanyID"`
}

func (u *PayRollReportItem) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}
func (m *PayRollReportItem) CheckPayment() {
	payRoll := PayRoll{}
	database.DB.Select("id").Find(&payRoll, "pay_roll_report_item_id = ?", m.ID)
	paymentTrans := Transaction{}
	count := int64(0)
	database.DB.Model(&paymentTrans).Where(&paymentTrans, "pay_roll_id = ? and is_pay_roll_payment = 1 and is_take_home_pay = 1", payRoll.ID).Count(&count)
	if m.Status != "FINISHED" {
		m.Status = "PAID"
	}
}

func (m *PayRollReportItem) MapToResp() resp.PayRollReportItemResponse {
	return resp.PayRollReportItemResponse{
		ID:                     m.ID,
		TotalTakeHomePay:       m.TotalTakeHomePay,
		TotalReimbursement:     m.TotalReimbursement,
		EmployeePhone:          m.EmployeePhone,
		EmployeeEmail:          m.EmployeeEmail,
		EmployeeName:           m.EmployeeName,
		EmployeeID:             m.EmployeeID,
		EmployeeIdentityNumber: m.EmployeeIdentityNumber,
		BankAccountNumber:      m.BankAccountNumber,
		BankName:               m.BankName,
		BankCode:               m.BankCode,
		BankID:                 m.BankID,
		Status:                 m.Status,
		PayRoll:                m.PayRoll.MapWithOutDetail(),
	}
}

func (m PayRollReportItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.MapToResp())
}
