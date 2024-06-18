package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayRollItem struct {
	Base
	ItemType         string        `json:"item_type" gorm:"type:enum('SALARY', 'ALLOWANCE', 'OVERTIME', 'DEDUCTION', 'REIMBURSEMENT')"`
	AccountPayableID *string       `json:"account_payable_id"`
	Title            string        `json:"title"`
	Notes            string        `json:"notes"`
	IsDefault        bool          `json:"is_default"`
	IsDeductible     bool          `json:"is_deductible"`
	IsTax            bool          `json:"is_tax"`
	TaxAutoCount     bool          `json:"tax_auto_count"`
	IsTaxCost        bool          `json:"is_tax_cost"`
	IsTaxAllowance   bool          `json:"is_tax_allowance"`
	Amount           float64       `json:"amount"`
	PayRollID        string        `json:"pay_roll_id"`
	PayRoll          PayRoll       `gorm:"foreignKey:PayRollID" json:"-"`
	ReimbursementID  *string       `json:"reimbursement_id"`
	Reimbursement    Reimbursement `gorm:"foreignKey:ReimbursementID" json:"-"`
	Bpjs             bool          `json:"bpjs"`
	BpjsCounted      bool          `json:"bpjs_counted"`
	Tariff           float64       `json:"tariff"`
	CompanyID        string        `json:"company_id" gorm:"not null"`
	Company          Company       `gorm:"foreignKey:CompanyID"`
}

type PayRollItemReq struct {
	ItemType        string  `json:"item_type"`
	Title           string  `json:"title"`
	Notes           string  `json:"notes"`
	IsDefault       bool    `json:"is_default"`
	IsDeductible    bool    `json:"is_deductible"`
	IsTax           bool    `json:"is_tax"`
	TaxAutoCount    bool    `json:"tax_auto_count"`
	IsTaxCost       bool    `json:"is_tax_cost"`
	IsTaxAllowance  bool    `json:"is_tax_allowance"`
	Amount          float64 `json:"amount"`
	PayRollID       string  `json:"pay_roll_id"`
	ReimbursementID *string `json:"reimbursement_id"`
}

func (u *PayRollItem) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m PayRollItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.PayRollItemResponse{})
}
