package model

import (
	"avolta/object/resp"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayRollItem struct {
	Base
	ItemType         string         `json:"item_type" gorm:"type:enum('SALARY', 'ALLOWANCE', 'OVERTIME', 'DEDUCTION', 'REIMBURSEMENT')"`
	AccountPayableID sql.NullString `json:"account_payable_id"`
	Title            string         `json:"title" binding:"required"`
	Notes            string         `json:"notes"`
	IsDefault        bool           `json:"is_default"`
	IsDeductible     bool           `json:"is_deductible"`
	IsTax            bool           `json:"is_tax"`
	TaxAutoCount     bool           `json:"tax_auto_count"`
	IsTaxCost        bool           `json:"is_tax_cost"`
	IsTaxAllowance   bool           `json:"is_tax_allowance"`
	Amount           float64        `json:"amount"`
	PayRollID        sql.NullString `json:"pay_roll_id" binding:"required"`
}

func (u *PayRollItem) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m PayRollItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.PayRollItemReponse{})
}
