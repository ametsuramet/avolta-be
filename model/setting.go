package model

import (
	"avolta/object/resp"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setting struct {
	Base
	PayRollAutoNumber                  bool    `json:"pay_roll_auto_number" gorm:"default:true"`
	PayRollAutoFormat                  string  `json:"pay_roll_auto_format" gorm:"default:'{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}'"`
	PayRollStaticCharacter             string  `json:"pay_roll_static_character" gorm:"default:'PAYROLL'"`
	PayRollAutoNumberCharacterLength   int     `json:"pay_roll_auto_number_character_length" gorm:"default:5"`
	PayRollPayableAccountID            *string `json:"pay_roll_payable_account_id"`
	PayRollPayableAccount              Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PayRollPayableAccountID" json:"pay_roll_payable_account"`
	PayRollExpenseAccountID            *string `json:"pay_roll_expense_account_id"`
	PayRollExpenseAccount              Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PayRollExpenseAccountID" json:"pay_roll_expense_account"`
	PayRollAssetAccountID              *string `json:"pay_roll_asset_account_id"`
	PayRollAssetAccount                Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PayRollAssetAccountID" json:"pay_roll_asset_account"`
	PayRollTaxAccountID                *string `json:"pay_roll_tax_account_id"`
	PayRollTaxAccount                  Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PayRollTaxAccountID" json:"pay_roll_tax_account"`
	PayRollCostAccountID               *string `json:"pay_roll_cost_account_id"`
	PayRollCostAccount                 Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PayRollCostAccountID" json:"pay_roll_cost_account"`
	ReimbursementPayableAccountID      *string `json:"reimbursement_payable_account_id"`
	ReimbursementPayableAccount        Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ReimbursementPayableAccountID" json:"reimbursement_payable_account"`
	ReimbursementExpenseAccountID      *string `json:"reimbursement_expense_account_id"`
	ReimbursementExpenseAccount        Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ReimbursementExpenseAccountID" json:"reimbursement_expense_account"`
	ReimbursementAssetAccountID        *string `json:"reimbursement_asset_account_id"`
	ReimbursementAssetAccount          Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ReimbursementAssetAccountID" json:"reimbursement_asset_account"`
	IsEffectiveRateAverage             bool    `json:"is_effective_rate_average"`
	IsGrossUp                          bool    `json:"is_gross_up"`
	BpjsKes                            bool    `json:"bpjs_kes"`
	BpjsTkJht                          bool    `json:"bpjs_tk_jht"`
	BpjsTkJkm                          bool    `json:"bpjs_tk_jkm"`
	BpjsTkJp                           bool    `json:"bpjs_tk_jp"`
	BpjsTkJkk                          bool    `json:"bpjs_tk_jkk"`
	IncentiveAutoNumber                bool    `json:"incentive_auto_number" gorm:"default:true"`
	IncentiveAutoFormat                string  `json:"incentive_auto_format" gorm:"default:'{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}'"`
	IncentiveStaticCharacter           string  `json:"incentive_static_character" gorm:"default:'INS'"`
	IncentiveAutoNumberCharacterLength int     `json:"incentive_auto_number_character_length" gorm:"default:5"`
}

func (u *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Setting) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.SettingResponse{
		ID:                                 m.ID,
		PayRollAutoNumber:                  m.PayRollAutoNumber,
		PayRollAutoFormat:                  m.PayRollAutoFormat,
		PayRollStaticCharacter:             m.PayRollStaticCharacter,
		PayRollAutoNumberCharacterLength:   m.PayRollAutoNumberCharacterLength,
		PayRollPayableAccountID:            m.PayRollPayableAccountID,
		PayRollExpenseAccountID:            m.PayRollExpenseAccountID,
		PayRollAssetAccountID:              m.PayRollAssetAccountID,
		PayRollTaxAccountID:                m.PayRollTaxAccountID,
		PayRollCostAccountID:               m.PayRollCostAccountID,
		ReimbursementPayableAccountID:      m.ReimbursementPayableAccountID,
		ReimbursementExpenseAccountID:      m.ReimbursementExpenseAccountID,
		ReimbursementAssetAccountID:        m.ReimbursementAssetAccountID,
		IsEffectiveRateAverage:             m.IsEffectiveRateAverage,
		IsGrossUp:                          m.IsGrossUp,
		BpjsKes:                            m.BpjsKes,
		BpjsTkJht:                          m.BpjsTkJht,
		BpjsTkJkm:                          m.BpjsTkJkm,
		BpjsTkJp:                           m.BpjsTkJp,
		BpjsTkJkk:                          m.BpjsTkJkk,
		IncentiveAutoNumber:                m.IncentiveAutoNumber,
		IncentiveAutoFormat:                m.IncentiveAutoFormat,
		IncentiveStaticCharacter:           m.IncentiveStaticCharacter,
		IncentiveAutoNumberCharacterLength: m.IncentiveAutoNumberCharacterLength,
	})
}
