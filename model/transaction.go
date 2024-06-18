package model

import (
	"avolta/object/resp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	Base
	Description            string        `json:"description"`
	Notes                  string        `json:"notes"`
	Credit                 float64       `json:"credit"`
	Debit                  float64       `json:"debit"`
	Amount                 float64       `json:"amount"`
	Date                   time.Time     `json:"date"`
	IsIncome               bool          `json:"is_income"`
	IsExpense              bool          `json:"is_expense"`
	IsJournal              bool          `json:"is_journal"`
	IsAccountReceivable    bool          `json:"is_account_receivable"`
	IsAccountPayable       bool          `json:"is_account_payable"`
	AccountSourceID        *string       `json:"account_source_id"`
	AccountSourceName      string        `json:"account_source_name" gorm:"-"`
	AccountSource          Account       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AccountSourceID" json:"-"`
	AccountDestinationID   *string       `json:"account_destination_id"`
	AccountDestinationName string        `json:"account_destination_name" gorm:"-"`
	AccountDestination     Account       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:AccountDestinationID" json:"-"`
	EmployeeID             string        `json:"employee_id"`
	Employee               Employee      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:EmployeeID" json:"-"`
	Images                 []Image       `json:"images" gorm:"-"`
	PayRollID              *string       `json:"pay_roll_id"`
	PayRoll                PayRoll       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PayRollID" json:"-"`
	ReimbursementID        *string       `json:"reimbursement_id"`
	Reimbursement          Reimbursement `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ReimbursementID" json:"-"`
	TaxPaymentID           string        `json:"tax_payment_id"`
	PayRollPayableID       string        `json:"pay_roll_payable_id"`
	IsPayRollPayment       bool          `json:"is_pay_roll_payment"`
	IsReimbursementPayment bool          `json:"is_reimbursement_payment"`
	TransactionRefID       *string       `json:"transaction_ref_id"`
	TransactionRefs        []Transaction `json:"transaction_refs" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:TransactionRefID"`
	IsTakeHomePay          bool          `json:"is_take_home_pay"`
	CompanyID              string        `json:"company_id"`
	Company                Company       `gorm:"foreignKey:CompanyID"`
}

func (u *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Transaction) MarshalJSON() ([]byte, error) {
	accountSourceId := ""
	accountSourceName := ""
	if m.AccountSourceID != nil {
		accountSourceId = *m.AccountSourceID
		accountSourceName = m.AccountSource.Name
	}
	accountDestinationId := ""
	accountDestinationName := ""
	if m.AccountDestinationID != nil {
		accountDestinationId = *m.AccountDestinationID
		accountDestinationName = m.AccountDestination.Name
	}

	transaction_refs := []resp.TransactionResponse{}
	for _, v := range m.TransactionRefs {
		transaction_refs = append(transaction_refs, resp.TransactionResponse{
			ID:                     v.ID,
			Description:            v.Description,
			Notes:                  v.Notes,
			Credit:                 v.Credit,
			Debit:                  v.Debit,
			Amount:                 v.Amount,
			Date:                   v.Date.Format("2006-01-02 15:04:05"),
			IsIncome:               v.IsIncome,
			IsExpense:              v.IsExpense,
			IsJournal:              v.IsJournal,
			IsAccountReceivable:    v.IsAccountReceivable,
			IsAccountPayable:       v.IsAccountPayable,
			AccountDestinationID:   *v.AccountDestinationID,
			AccountDestinationName: v.AccountDestination.Name,
			EmployeeID:             v.EmployeeID,
			EmployeeName:           v.Employee.FullName,
		})
	}
	return json.Marshal(resp.TransactionResponse{
		ID:                     m.ID,
		Description:            m.Description,
		Notes:                  m.Notes,
		Credit:                 m.Credit,
		Debit:                  m.Debit,
		Amount:                 m.Amount,
		Date:                   m.Date.Format("2006-01-02 15:04:05"),
		IsIncome:               m.IsIncome,
		IsExpense:              m.IsExpense,
		IsJournal:              m.IsJournal,
		IsAccountReceivable:    m.IsAccountReceivable,
		IsAccountPayable:       m.IsAccountPayable,
		AccountSourceID:        accountSourceId,
		AccountSourceName:      accountSourceName,
		AccountDestinationID:   accountDestinationId,
		AccountDestinationName: accountDestinationName,
		EmployeeID:             m.EmployeeID,
		EmployeeName:           m.Employee.FullName,
		TransactionRefs:        transaction_refs,
		IsPayRollPayment:       m.IsPayRollPayment,
		IsReimbursementPayment: m.IsReimbursementPayment,
	})
}
