package model

import (
	"avolta/object/resp"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	Base
	Name                  string        `json:"name" bson:"name"`
	Code                  string        `json:"code" bson:"code"`
	Color                 string        `json:"color" bson:"color"`
	Description           string        `json:"description" bson:"description"`
	IsDeletable           bool          `json:"is_deletable" bson:"is_deletable"`
	IsReport              bool          `json:"is_report" bson:"is_report" gorm:"-"`
	IsAccountReport       bool          `json:"is_account_report" bson:"is_account_report" gorm:"-"`
	IsCashflowReport      bool          `json:"is_cashflow_report" bson:"is_cashflow_report" gorm:"-"`
	IsPDF                 bool          `json:"is_pdf" bson:"is_pdf" gorm:"-"`
	Type                  string        `json:"type" bson:"type"`
	Category              string        `json:"category" bson:"category"`
	CashflowGroup         string        `json:"cashflow_group" bson:"cashflow_group"`
	CashflowSubGroup      string        `json:"cashflow_subgroup" bson:"cashflow_group"`
	Transactions          []Transaction `gorm:"-"`
	IsTax                 bool          `json:"is_tax" bson:"is_tax" gorm:"default:false"`
	TypeLabel             string        `gorm:"-" json:"type_label"`
	CashflowGroupLabel    string        `gorm:"-" json:"cashflow_group_label"`
	CashflowSubGroupLabel string        `gorm:"-" json:"cashflow_subgroup_label" `
	CompanyID             string        `json:"company_id"`
	Company               Company       `gorm:"foreignKey:CompanyID"`
}

func (u *Account) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	fmt.Println(ctx)
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}

	// if u.CompanyID == "" {
	// 	tx.Statement.SetColumn("company_id", ctx.Value("ID-Company"))
	// }
	return
}

func (m Account) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.AccountResponse{
		ID:                    m.ID,
		Name:                  m.Name,
		Code:                  m.Code,
		Color:                 m.Color,
		Description:           m.Description,
		IsDeletable:           m.IsDeletable,
		IsReport:              m.IsReport,
		IsAccountReport:       m.IsAccountReport,
		IsCashflowReport:      m.IsCashflowReport,
		IsPDF:                 m.IsPDF,
		Type:                  m.Type,
		Category:              m.Category,
		CashflowGroup:         m.CashflowGroup,
		CashflowSubGroup:      m.CashflowSubGroup,
		IsTax:                 m.IsTax,
		TypeLabel:             m.TypeLabel,
		CashflowGroupLabel:    m.CashflowGroupLabel,
		CashflowSubGroupLabel: m.CashflowSubGroupLabel,
	})
}
