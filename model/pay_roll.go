package model

import (
	"avolta/object/resp"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayRoll struct {
	Base
	Title                  string         `json:"title"`
	Notes                  string         `json:"notes"`
	StartDate              time.Time      `json:"start_date" binding:"required"`
	EndDate                time.Time      `json:"end_date" binding:"required"`
	Files                  string         `json:"files" gorm:"default:'[]'"`
	TotalIncome            float64        `json:"total_income"`
	TotalReimbursement     float64        `json:"total_reimbursement"`
	TotalDeduction         float64        `json:"total_deduction"`
	TotalTax               float64        `json:"total_tax"`
	TaxCost                float64        `json:"tax_cost"`
	NetIncome              float64        `json:"net_income"`
	NetIncomeBeforeTaxCost float64        `json:"net_income_before_tax_cost"`
	TakeHomePay            float64        `json:"take_home_pay"`
	TotalPayable           float64        `json:"total_payable"`
	TaxAllowance           float64        `json:"tax_allowance"`
	IsGrossUp              bool           `json:"is_gross_up"`
	IsEffectiveRateAverage bool           `json:"is_effective_rate_average"`
	Status                 string         `json:"status" gorm:"type:enum('DRAFT', 'RUNNING', 'FINISHED');default:'DRAFT'"`
	Attachments            []string       `json:"attachments" gorm:"-"`
	Transactions           []Transaction  `json:"transactions" gorm:"-"`
	PayableTransactions    []Transaction  `json:"payable_transactions" gorm:"-"`
	Items                  []PayRollItem  `json:"items" gorm:"-"`
	TakeHomePayCounted     string         `json:"take_home_pay_counted" gorm:"-"`
	TaxPaymentID           sql.NullString `json:"tax_payment_id"`
	EmployeeID             string         `binding:"required" json:"employee_id"`
	Employee               Employee       `gorm:"foreignKey:EmployeeID" `
}

func (u *PayRoll) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m PayRoll) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.PayRollReponse{
		ID:           m.ID,
		Notes:        m.Notes,
		StartDate:    m.StartDate.Format("2006-01-02"),
		EndDate:      m.EndDate.Format("2006-01-02"),
		EmployeeID:   m.EmployeeID,
		EmployeeName: m.Employee.FullName,
	})
}
