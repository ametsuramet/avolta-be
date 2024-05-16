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
	Description          string    `json:"description"`
	Notes                string    `json:"notes"`
	Credit               float64   `json:"credit"`
	Debit                float64   `json:"debit"`
	Amount               float64   `json:"amount"`
	Date                 time.Time `json:"date"`
	IsIncome             bool      `json:"is_income"`
	IsExpense            bool      `json:"is_expense"`
	IsJournal            bool      `json:"is_journal"`
	IsAccountReceivable  bool      `json:"is_account_receivable"`
	IsAccountPayable     bool      `json:"is_account_payable"`
	AccountSourceID      string    `json:"account_source_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	EmployeeID           string    `json:"employee_id"`
	Images               []Image   `json:"images" gorm:"-"`
	PayRollID            string    `json:"pay_roll_id"`
}

func (u *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.TransactionReponse{})
}
